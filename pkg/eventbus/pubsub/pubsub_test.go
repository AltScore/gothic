//go:build integration

package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/AltScore/gothic/pkg/eventbus"
	"github.com/AltScore/gothic/test/pubsubtest"
	"github.com/modernice/goes/codec"
	"github.com/modernice/goes/event"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"math/rand"

	"strconv"
	"testing"
	"time"
)

const (
	// Increase this to debug
	testDeadlineSeconds = 120
	testEventName       = "test-event"
)

type PubSubTestSuite struct {
	suite.Suite

	run   int
	topic string

	pubsubtest.LocalEmulator
	registry *codec.Registry
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
	s.run = rand.Int()
	s.topic = "test-topic." + strconv.Itoa(s.run)

	s.registry = codec.New(codec.Debug(true))

	// PublisherAdapter requires a codec that can encode the event
	codec.Register[string](s.registry, testEventName)
}

func (s *PubSubTestSuite) Test_List_topics() {
	ctx, cancelCtx := context.WithTimeout(context.Background(), testDeadlineSeconds*time.Second)
	defer cancelCtx()

	it := s.Client().Topics(ctx)

	var topics []*pubsub.Topic

	var err error
	for {
		var topic *pubsub.Topic
		topic, err = it.Next()
		if err != nil {
			break
		}
		topics = append(topics, topic)
	}

	s.Require().Equal(iterator.Done, err)

	s.T().Logf("Topics: %v", topics)
}

func (s *PubSubTestSuite) Test_Given_a_pubsub_publisher_When_publish_Then_message_was_published() {
	ctx, cancelCtx := context.WithTimeout(context.Background(), testDeadlineSeconds*time.Second)
	defer cancelCtx()

	// GIVEN a topic and a subscription
	topic := s.givenTopic(ctx)

	subs := s.givenSubscription(topic, ctx)

	// AND a publisher
	publisher := s.givenAPublisher(topic, ctx)

	// WHEN we send a message to the topic
	ev := event.New(testEventName, "sample data")

	err := publisher.Publish(ev.Any())

	// THEN no error should be produced
	s.Require().NoError(err)

	fmt.Printf("Published message %s\n", ev.ID())

	// THEN we should receive the message on the subscription
	msgReceived := ""

	go func() {
		fmt.Printf("Waiting for message on subscription %s\n", subs.ID())
		err = subs.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			msgReceived = string(msg.Data)
			fmt.Printf("Got message: %s\n", msgReceived)
			cancelCtx()
		})
	}()

	fmt.Printf("Waiting for subscription to finish\n")
	<-ctx.Done()
	fmt.Printf("Context done: %s\n", ctx.Err())

	fmt.Printf("End waiting for message on subscription %s\n", subs.ID())
	s.Require().NoError(err)

	s.Require().Equal("\"sample data\"", msgReceived)
}

func (s *PubSubTestSuite) givenAPublisher(topic *pubsub.Topic, ctx context.Context) *Publisher {
	config := PublisherConfig{
		ProjectID: s.ProjectID(),
		TopicName: topic.ID(),
	}

	publisher := NewPublisher(ctx, s.Client(), s.registry, zap.NewExample(), config)
	return publisher
}

func (s *PubSubTestSuite) Test_can_send_and_receive() {
	ctx, cancelCtx := context.WithTimeout(context.Background(), testDeadlineSeconds*time.Second)
	defer cancelCtx()

	// GIVEN a topic and a subscription
	topic := s.givenTopic(ctx)
	sub := s.givenSubscription(topic, ctx)

	// AND an adapted localBus
	localBus := newLocalBusMock()

	options := PullAdapterConfig{
		ProjectID:        s.ProjectID(),
		SubscriptionName: sub.ID(),
		Debug:            true,
	}

	pullAdapter := NewPullAdapter(s.Client(), localBus, s.registry, zap.NewExample(), options)

	localBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

	// WHEN we send a message to the topic
	go func() {
		fmt.Printf("Starting pull adapter\n")
		err := pullAdapter.Start(ctx)
		s.Require().NoError(err)
	}()

	publisher := s.givenAPublisher(topic, ctx)
	ev := event.New("test-event", "sample data")

	err := publisher.Publish(ev.Any())
	fmt.Printf("Published message %s\n", ev.ID())

	// THEN no error should be produced
	s.Require().NoError(err)

	// THEN we should receive the message on the subscription
	select {
	case <-ctx.Done():
		s.Fail("Context done")
	case <-localBus.done:
		fmt.Printf("Message received and local bus done\n")
		localBus.AssertCalled(s.T(), "Publish", mock.Anything, mock.Anything)
	}
}

func (s *PubSubTestSuite) givenSubscription(topic *pubsub.Topic, ctx context.Context) *pubsub.Subscription {
	cfg := pubsub.SubscriptionConfig{
		Topic: topic,
	}

	subs, err := s.Client().CreateSubscription(ctx, s.topic+"-subscription-1", cfg)
	s.Require().NoError(err)

	fmt.Printf("Created subscription %s\n", subs.ID())

	return subs
}

func (s *PubSubTestSuite) givenTopic(ctx context.Context) *pubsub.Topic {
	topic, err := s.Client().CreateTopic(ctx, s.topic)
	s.Require().NoError(err, "Could not create topic %s\n", s.topic)

	fmt.Printf("Created topic %s", topic.ID())

	s.Require().NoError(err)

	return topic
}

type localBusMock struct {
	mock.Mock

	done chan struct{}
}

func newLocalBusMock() *localBusMock {
	return &localBusMock{
		done: make(chan struct{}),
	}
}

func (p *localBusMock) Publish(event eventbus.Event, options ...eventbus.Option) error {
	defer func() {
		fmt.Printf("Local bus done\n")
		p.done <- struct{}{}
	}()

	fmt.Printf("Received event %s\n", event.ID())

	args := p.Called(event, options)

	return args.Error(0)
}
