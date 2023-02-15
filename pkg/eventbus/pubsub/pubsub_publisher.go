package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/AltScore/gothic/pkg/errors"
	"github.com/AltScore/gothic/pkg/eventbus"
	"github.com/AltScore/gothic/pkg/logger"
	"github.com/modernice/goes/codec"
	"github.com/modernice/goes/event"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type PublisherConfig struct {
	ProjectID  string
	TopicName  string
	LogMessage bool
}

type Publisher struct {
	ctx    context.Context
	client *pubsub.Client

	topic    *pubsub.Topic
	logger   logger.Logger
	config   PublisherConfig
	encoding codec.Encoding
}

// NewPublisher creates a new Publisher that publishes events to a PubSub topic
// The topic must be created before using this gateway
// Messages are sent in order, so the OrderingKey is set to the Aggregate ID
//
// To authenticate with PubSub, the GOOGLE_APPLICATION_CREDENTIALS environment variable must be set
// See https://cloud.google.com/docs/authentication/getting-started for more information
func NewPublisher(ctx context.Context, client *pubsub.Client, encoding codec.Encoding, log logger.Logger, config PublisherConfig) *Publisher {
	errors.EnsureNotNil(client, "PubSub client")
	errors.EnsureNotNil(encoding, "Event encoding")

	log.Info("Connected to PubSub", zap.String("project_id", config.ProjectID), zap.String("topic_name", config.TopicName))

	topic := client.Topic(config.TopicName)
	topic.EnableMessageOrdering = true // This is required for the OrderingKey to work. It is critical for Aggregate Event Sourcing

	return &Publisher{
		ctx:      ctx,
		client:   client,
		encoding: encoding,
		topic:    topic,
		logger:   log,
		config:   config,
	}
}

// Publish sends the given events to the configured PubSub topic
// Each message is sent in order, if an error is produced, it stops sending and returns the error.
//
// The event name can be found in the "type" attribute of the message
func (g *Publisher) Publish(e event.Event, options ...eventbus.Option) error {
	envelope := &eventbus.EventEnvelope{
		Event: e,
		Ctx:   g.ctx,
	}

	envelope.ProcessOptions(options)

	data, err := g.encoding.Marshal(e.Data())
	if err != nil {
		return err
	}

	aggID, aggName, aggVersion := e.Aggregate()

	start := time.Now()
	result := g.topic.Publish(envelope.Ctx, &pubsub.Message{
		Data:        data,
		OrderingKey: aggID.String(),
		Attributes: map[string]string{
			EventIDMessageAttributeKey:          e.ID().String(),
			EventNameMessageAttributeKey:        e.Name(),
			EventTimeMessageAttributeKey:        e.Time().Format(EventTimeFormat),
			AggregateIDMessageAttributeKey:      aggID.String(),
			AggregateNameMessageAttributeKey:    aggName,
			AggregateVersionMessageAttributeKey: strconv.Itoa(aggVersion),
		},
	})

	_, err = result.Get(envelope.Ctx)
	// TODO recover in case of errors with ordered messages
	if err != nil {
		return err
	}

	if g.config.LogMessage {
		end := time.Now()
		g.logger.Info("Message sent", zap.String("type", envelope.Event.Name()), zap.Any("aggregateID", aggID), zap.Duration("latency", end.Sub(start)))
	}

	return nil
}
