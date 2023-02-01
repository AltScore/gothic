package pubsub

import (
	"bitbucket.org/altscore/altscore-credits-api.git/pkg/app/errors"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/AltScore/gothic/pkg/es/event"
	"github.com/AltScore/gothic/pkg/eventbus"
	"github.com/AltScore/gothic/pkg/logger"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type PullAdapterConfig struct {
	ProjectID        string `yaml:"project_id"`
	SubscriptionName string `yaml:"subscription_name"`
}

// PullAdapter pulls events from a PubSub topic and publishes them to a local event bus
type PullAdapter struct {
	client    *pubsub.Client
	logger    logger.Logger
	config    PullAdapterConfig
	publisher eventbus.Publisher
	sub       *pubsub.Subscription
}

// NewPullAdapter creates a new PullAdapter that pulls events from a PubSub topic and publishes them to another publisher
// To authenticate with PubSub, the GOOGLE_APPLICATION_CREDENTIALS environment variable must be set
// See https://cloud.google.com/docs/authentication/getting-started for more information
func NewPullAdapter(client *pubsub.Client, publisher eventbus.Publisher, log logger.Logger, config PullAdapterConfig) *PullAdapter {
	errors.EnsureNotNil(client, "client")

	log.Info("Connected to PubSub", zap.String("project_id", config.ProjectID), zap.String("subscription", config.SubscriptionName))

	sub := client.Subscription(config.SubscriptionName)

	return &PullAdapter{
		client:    client,
		logger:    log,
		config:    config,
		publisher: publisher,
		sub:       sub,
	}
}

// Start starts the adapter
func (a *PullAdapter) Start(ctx context.Context) error {
	err := a.sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		if err := a.publish(ctx, m); err != nil {
			m.Nack()
		} else {
			m.Ack() // Acknowledge that we've consumed the message.
		}
	})
	return err
}

func (a *PullAdapter) publish(ctx context.Context, m *pubsub.Message) error {
	ev, err := a.unmarshalEvent(m)
	if err != nil {
		return err
	}

	a.logger.Info("Received event from PubSub", zap.String("event", ev.Name()), zap.Any("id", ev.ID()), zap.Any("agg_id", ev.AggregateID()))

	errCh := make(chan error)

	err = a.publisher.Publish(ev, eventbus.WithAckChan(errCh), eventbus.WithContext(ctx))

	if err != nil {
		return err
	}

	return <-errCh // Wait for the response to be received
}

func (a *PullAdapter) unmarshalEvent(msg *pubsub.Message) (*event.Event, error) {
	evID := msg.Attributes[EventIDMessageAttributeKey]
	aggID := msg.Attributes[AggregateIDMessageAttributeKey]
	aggName := msg.Attributes[AggregateNameMessageAttributeKey]
	aggVersionStr := msg.Attributes[AggregateVersionMessageAttributeKey]
	aggVersion, err := strconv.Atoi(aggVersionStr)

	if err != nil {
		return nil, fmt.Errorf("failed to parse aggregate version '%s': %w", aggVersionStr, err)
	}

	evTime, err := time.Parse(EventTimeFormat, msg.Attributes[EventTimeMessageAttributeKey])

	if err != nil {
		return nil, fmt.Errorf("failed to parse event time '%s': %w", msg.Attributes[EventTimeMessageAttributeKey], err)
	}

	eventName := msg.Attributes[EventNameMessageAttributeKey]

	data := map[string]interface{}{}

	ev := event.New(
		eventName,
		data,
		event.WithTypeAndVersion(aggName, aggID, aggVersion),
		event.WithID(event.ID(evID)),
		event.WithTime(evTime),
	)

	return &ev, nil
}
