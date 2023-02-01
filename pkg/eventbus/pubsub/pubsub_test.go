//go:build integration

package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/AltScore/gothic/pkg/es/event"
	"github.com/AltScore/gothic/test/pubsubtest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"strconv"
	"testing"
	"time"
)

type PubSubTestSuite struct {
	suite.Suite

	run   int
	topic string

	pubsubtest.LocalEmulator
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPubSubTestSuite(t *testing.T) {
	suite.Run(t, &PubSubTestSuite{})
}

func (s *PubSubTestSuite) SetupSuite() {
	s.LocalEmulator.Connect()
}

func (s *PubSubTestSuite) TearDownSuite() {
	s.LocalEmulator.Disconnect()
}

func (s *PubSubTestSuite) SetupTest() {
	s.run++
	s.topic = "test-topic-" + strconv.Itoa(s.run)
}

func (s *PubSubTestSuite) Test_Given_a_pubsub_publisher_When_publish_Then_message_was_published() {
	ctx := context.Background()

	// GIVEN a topic and a subscription
	ctx2, cancel1 := context.WithTimeout(ctx, 10*time.Second)
	defer cancel1()

	topic, err := s.Client().CreateTopic(ctx2, s.topic)
	s.Require().NoError(err)

	cfg := pubsub.SubscriptionConfig{
		Topic: topic,
	}

	subs, err := s.Client().CreateSubscription(ctx, s.topic+"-subscription-1", cfg)
	s.Require().NoError(err)

	// AND a publisher
	config := PublisherConfig{
		ProjectID: s.ProjectID(),
		TopicName: topic.ID(),
	}
	publisher := NewPublisher(ctx, s.Client(), zap.NewExample(), config)

	// WHEN we send a message to the topic
	ev := event.New("test-event", "sample data")

	err = publisher.Publish(ev)

	// THEN no error should be produced
	s.Require().NoError(err)

	// THEN we should receive the message on the subscription
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)

	msgReceived := ""

	err = subs.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msgReceived = string(msg.Data)
		cancel()
	})

	s.Require().NoError(err)

	s.Require().Equal("sample data", msgReceived)
}

func (s *PubSubTestSuite) Test_can_send_and_receive() {
	// GIVEN a topic and a subscription

	// WHEN we send a message to the topic

	// THEN we should receive the message on the subscription

}
