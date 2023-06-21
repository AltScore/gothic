package xmongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ReplicaSetName = "rs0"
)

type Option func(*MongoInMemory)

// MongoInMemory allows to start a MongoDB instance in memory to run tests against
// To use it, create a new instance of MongoInMemory and call Connect() before running your tests
// and call Disconnect() after your tests are done.
// The client to connect to MongoDB can be retrieved by calling Client().
type MongoInMemory struct {
	dbClient      *mongo.Client
	pool          *dockertest.Pool
	resource      *dockertest.Resource
	mongoPort     string
	container     *Container
	useReplicaSet bool
	logContainer  bool
	debug         bool
	containerName string
}

// Connect starts a MongoDB instance in memory and connects to it.
// You can pass options to configure the instance.
// The instance will be stopped and removed when Disconnect() is called.
// If you want to see the logs of the instance, pass the WithLogContainer() option.
// If you want to use a replica set, pass the WithReplicaSet() option.
// If you want to see the more logs, pass the WithDebug() option.
func (m *MongoInMemory) Connect(options ...Option) {
	for _, option := range options {
		option(m)
	}

	var err error
	m.pool, err = dockertest.NewPool("")

	m.failOnError(err, "Could not construct pool: %s")

	m.pool.MaxWait = 10 * time.Second

	m.failOnError(m.pool.Client.Ping(), "Could not connect to Docker: %s")

	m.containerName = "mongo-" + strings.ReplaceAll(time.Now().Format("15:04:05"), ":", "")

	envs := []string{
		"MONGODB_ADVERTISED_HOSTNAME=localhost",
		"ALLOW_EMPTY_PASSWORD=yes",
	}

	if m.useReplicaSet {
		envs = append(envs,
			"MONGODB_REPLICA_SET_MODE=primary",
			"MONGODB_REPLICA_SET_NAME="+ReplicaSetName,
		)
	}

	m.resource, err = m.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "bitnami/mongodb",
		Tag:        "4.4",
		Name:       m.containerName,
		Env:        envs,
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017/tcp": {
				{
					HostIP:   "0.0.0.0",
					HostPort: "27017",
				},
			},
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	m.failOnError(err, "Could not start resource: %s")

	if m.logContainer {
		m.container = &Container{
			pool:     m.pool,
			resource: m.resource,
		}

		go func() {
			m.failOnError(m.container.TailLogs(context.Background(), os.Stdout, true), "Could not tail logs: %v")
		}()
		log.Println("Showing logs for container: ", m.container.resource.Container.Name)
	}

	m.mongoPort = m.resource.GetPort("27017/tcp")

	log.Println("MongoDB running on port: ", m.mongoPort, " - ", m.mongoPort)

	m.failOnError(m.connectToMongo(), "Could not connect to docker: %s")
	log.Println("Connected to docker")
}

func (m *MongoInMemory) connectToMongo() error {
	return m.connectAndRun(func() error {
		return m.dbClient.Ping(context.TODO(), nil)
	})
}

func (m *MongoInMemory) connectAndRun(f func() error) error {
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	return m.pool.Retry(func() error {
		var (
			err           error
			connectionStr string
		)

		var optionsStr string
		if m.useReplicaSet {
			optionsStr = "?replicaSet=" + ReplicaSetName
		} else {
			optionsStr = ""
		}
		connectionStr = fmt.Sprintf("mongodb://localhost:%s/%s", m.mongoPort, optionsStr)

		clientOptions := options.Client().ApplyURI(connectionStr)

		m.dbClient, err = mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			m.log("Could not connect to mongo: %s", err)
			return err
		}

		err = f()

		m.logError(err, "Could not execute command")

		return err
	})
}

// Disconnect stops and removes the container
func (m *MongoInMemory) Disconnect() {
	// disconnect mongodb client
	m.disconnectMongoClient()

	if m.container != nil {
		m.container.Close()
	}

	// When you're done, kill and remove the container
	if err := m.pool.Purge(m.resource); err != nil {
		log.Printf("Could not purge resource: %s\n", err)
	}

	m.resource = nil
	m.pool = nil
}

func (m *MongoInMemory) disconnectMongoClient() {
	if m.dbClient == nil {
		return
	}

	if err := m.dbClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	m.dbClient = nil
}

func (m *MongoInMemory) Client() *mongo.Client {
	return m.dbClient
}

// initiateReplicaSet initiates a replica set with a single node in the mongo container.
// This is required for transactions to work.
func (m *MongoInMemory) initiateReplicaSet() error {
	err := backoff.Retry(func() error {
		output, err := m.executeShellCommand("rs.initiate({_id: 'rs0', members: [{_id: 0, host: '" + m.containerName + ":27017'}]})")

		if err != nil {
			m.log("Error sending command to mongo: %v", err)
			return err
		}

		m.log("Command sent to mongo: %v", output)

		if !strings.Contains(output, `{ "ok" : 1 }`) {
			return errors.New("not yet started")
		}
		return nil
	}, newExponentialBackOff())

	if err != nil {
		m.log("Error initiating replica set: %v", err)
		return err
	}

	log.Println("Initiate Replica set command sent")

	return backoff.Retry(func() error {
		status, err := m.executeShellCommand("rs.status()")

		if err != nil {
			log.Println("Error initiating replica set: ", err)
			return err
		}

		log.Println("Replica set status: ", status)

		if !strings.Contains(status, "PRIMARY") {
			return errors.New("not yet initiated replica set")
		}
		log.Println("Replica set initiated")
		return nil
	}, newExponentialBackOff())
}

func (m *MongoInMemory) log(fmtStr string, args ...any) {
	if m.debug {
		log.Printf(fmtStr, args...)
	}
}

func newExponentialBackOff() *backoff.ExponentialBackOff {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     50 * time.Millisecond,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         10 * time.Second,
		MaxElapsedTime:      60 * time.Second,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	return b
}

// executeShellCommand executes a mongo shell command in the mongo container.
func (m *MongoInMemory) executeShellCommand(command string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := docker.CreateExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Context:      ctx,
		Container:    m.resource.Container.ID,
		Cmd: []string{
			"mongo",
			"-u",
			"root",
			"-p",
			"password",
			"admin",
			"--eval",
			command,
		},
	}
	exec, err := m.pool.Client.CreateExec(opts)

	if err != nil {
		return "", err
	}

	output := bytes.NewBuffer(nil)

	execOptions := docker.StartExecOptions{
		Context:      ctx,
		OutputStream: output,
		ErrorStream:  os.Stderr,
		Detach:       false,
		Tty:          false,
		RawTerminal:  false,
		Success:      nil,
	}

	if m.debug {
		log.Println("Executing shell command: ", command)
	}
	err = m.pool.Client.StartExec(exec.ID, execOptions)
	return output.String(), err
}

func (m *MongoInMemory) failOnError(err error, fmtStr string) {
	if err != nil {
		panic(fmt.Sprintf(fmtStr, err))
	}
}

func (m *MongoInMemory) logError(err error, s string) {
	if err != nil && m.debug {
		m.log("%s: %v", s, err)
	}
}

// WithReplicaSet enables replica set in the mongo database. This is useful when you want to test transactions.
func WithReplicaSet() Option {
	return func(mim *MongoInMemory) {
		mim.useReplicaSet = true
	}
}

// WithContainerLogs enables logging of the container stdout and stderr
func WithContainerLogs() Option {
	return func(mim *MongoInMemory) {
		mim.logContainer = true
	}
}

// WithDebug enables debug mode. More logs will be printed
func WithDebug() Option {
	return func(mim *MongoInMemory) {
		mim.debug = true
	}
}
