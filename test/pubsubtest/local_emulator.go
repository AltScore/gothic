package pubsubtest

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"math/rand"
	"os"
	"time"
)

// LocalEmulator allows to start a PubSub emulator in memory to run tests against
// To use it, create a new instance of LocalEmulator and call Connect() before running your tests
// and call Disconnect() after your tests are done.
// The client to connect to PubSub can be retrieved by calling Client().
type LocalEmulator struct {
	client   *pubsub.Client
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (m *LocalEmulator) Connect() {
	rand.Seed(time.Now().UnixNano())

	var err error
	log.Println("Connecting to Docker")
	m.pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = m.pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image for version 5.0
	log.Println("Starting local PubSub emulator")
	m.resource, err = m.pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "pubsub-emulator",
		Repository: "gcr.io/google.com/cloudsdktool/google-cloud-cli",
		Tag:        "emulators",
		// gcloud beta emulators pubsub start --project=credit-flow-staging --host-port=0.0.0.0:8085 --verbosity=debug --user-output-enabled=true --log-http
		Cmd: []string{"gcloud", "beta", "emulators", "pubsub", "start", "--project=" + m.ProjectID(),
			"--host-port=0.0.0.0:8085",
			"--verbosity=debug", "--user-output-enabled=true", "--log-http"},
		Env: []string{},
		ExposedPorts: []string{
			"8085",
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

	m.pool.MaxWait = 60 * time.Second

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	iteration := 0

	err = m.pool.Retry(func() error {
		log.Println("Trying to connect to pubsub container")
		var err error
		port := m.resource.GetPort("8085/tcp")

		log.Printf("Connecting to pubsub container on port '%s'", port)
		// Set PUBSUB_EMULATOR_HOST environment variable. https://wahlstrand.dev/posts/2021-07-11-testing-pubsub-locally/
		err = os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:"+port)
		if err != nil {
			log.Fatalf("Could not set PUBSUB_EMULATOR_HOST environment variable: %s", err)
		}
		m.client, err = pubsub.NewClient(context.TODO(), m.ProjectID())

		if err == nil {
			iteration++
			err = m.ping(iteration)
		}

		if err != nil {
			log.Printf("Could not connect to pubsub container: %s", err)
		} else {
			log.Printf("Connected to pubsub container on port '%s'", port)
		}

		return err
	})

	if err != nil {
		log.Printf("Could not connect to pubsub container: %s", err)
	}
}

func (m *LocalEmulator) ping(iteration int) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelCtx()

	pingTopic := fmt.Sprintf("ping.%d.%d", iteration, rand.Int())
	_, err := m.Client().CreateTopic(ctx, pingTopic)

	fmt.Printf("Ping Topic %s: %v", pingTopic, err)
	return err
}

func (m *LocalEmulator) ProjectID() string {
	return "test-project"
}

func (m *LocalEmulator) Disconnect() {
	log.Printf("Disconnecting from pubsub container")
	if err := m.client.Close(); err != nil {
		log.Printf("Could not close client: %v\n", err)
	}

	// When you're done, kill and remove the container
	log.Printf("Killing and removing pubsub container")
	if err := m.pool.Purge(m.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	log.Printf("Finished disconnecting from pubsub container")
	m.resource = nil
	m.pool = nil
	m.client = nil
}

func (m *LocalEmulator) Client() *pubsub.Client {
	return m.client
}
