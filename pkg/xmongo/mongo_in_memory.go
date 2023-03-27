package xmongo

import (
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// MongoInMemory allows to start a MongoDB instance in memory to run tests against
// To use it, create a new instance of MongoInMemory and call Connect() before running your tests
// and call Disconnect() after your tests are done.
// The client to connect to MongoDB can be retrieved by calling Client().
type MongoInMemory struct {
	dbClient *mongo.Client
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (m *MongoInMemory) Connect() {
	var err error
	m.pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = m.pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image for version 5.0
	m.resource, err = m.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			// username and password for mongodb superuser
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = m.pool.Retry(func() error {
		var err error
		connectionStr := fmt.Sprintf("mongodb://root:password@localhost:%s", m.resource.GetPort("27017/tcp"))
		m.dbClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(connectionStr),
		)
		if err != nil {
			return err
		}
		return m.dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}

func (m *MongoInMemory) Disconnect() {
	// When you're done, kill and remove the container
	if err := m.pool.Purge(m.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// disconnect mongodb client
	if err := m.dbClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	m.resource = nil
	m.pool = nil
	m.dbClient = nil
}

func (m *MongoInMemory) Client() *mongo.Client {
	return m.dbClient
}
