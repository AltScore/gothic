package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/AltScore/gothic/pkg/errors"
	"github.com/AltScore/gothic/pkg/eventbus"
	"github.com/AltScore/gothic/pkg/logger"
	"github.com/google/uuid"
	"github.com/modernice/goes/codec"
	"github.com/modernice/goes/event"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type PullAdapterConfig struct {
	ProjectID        string
	SubscriptionName string
	Debug            bool
}

// PullAdapter pulls events from a PubSub topic and publishes them to a local event bus
type PullAdapter struct {
	client    *pubsub.Client
	logger    logger.Logger
	config    PullAdapterConfig
	publisher eventbus.Publisher
	encoding  codec.Encoding
	sub       *pubsub.Subscription
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewPullAdapter creates a new PullAdapter that pulls events from a PubSub topic and publishes them to another publisher
//
// The provided enconder should have all the event types registered before processing events.
//
// To authenticate with PubSub, the GOOGLE_APPLICATION_CREDENTIALS environment variable must be set
// See https://cloud.google.com/docs/authentication/getting-started for more information
func NewPullAdapter(client *pubsub.Client, publisher eventbus.Publisher, encoder codec.Encoding, log logger.Logger, config PullAdapterConfig) *PullAdapter {
	errors.EnsureNotNil(client, "client")
	errors.EnsureNotNil(encoder, "encoder")
	errors.EnsureNotNil(encoder, "encoder")

	log.Info("Connected to PubSub", zap.String("project_id", config.ProjectID), zap.String("subscription", config.SubscriptionName))

	sub := client.Subscription(config.SubscriptionName)

	pa := &PullAdapter{
		client:    client,
		logger:    log,
		config:    config,
		publisher: publisher,
		encoding:  encoder,
		sub:       sub,
	}

	if lp, ok := publisher.(eventbus.LifeCycleProvider); ok {
		lp.AddLifecycleListener(lifeCycleListener{pa})

	}

	return pa
}

// Start starts the adapter
func (a *PullAdapter) start(ctx context.Context) error {
	a.ctx, a.cancel = context.WithCancel(ctx)

	go func() {
		err := a.sub.Receive(a.ctx, func(ctx context.Context, m *pubsub.Message) {
			if err := a.publish(ctx, m); err != nil {
				m.Nack()
			} else {
				m.Ack() // Acknowledge that we've consumed the message.
			}
		})
		if err != nil {
			a.logger.Error("Failed to receive from PubSub", zap.Error(err))
		}
	}()

	return nil
}

func (a *PullAdapter) stop() {
	if a.cancel != nil {
		a.logger.Info("Stopping PubSub adapter, canceling context")
		a.cancel()
		a.cancel = nil
		a.ctx = nil
	}
}

func (a *PullAdapter) publish(ctx context.Context, m *pubsub.Message) error {
	ev, err := a.unmarshalEvent(m)
	if err != nil {
		return err
	}

	aggID, _, _ := ev.Aggregate()

	a.logger.Info("Received event from PubSub", zap.String("event", ev.Name()), zap.Any("id", ev.ID()), zap.Any("agg_id", aggID))

	errCh := make(chan error)

	err = a.publisher.Publish(ev, eventbus.WithAckChan(errCh), eventbus.WithContext(ctx))

	if err != nil {
		return err
	}

	return <-errCh // Wait for the response to be received
}

func (a *PullAdapter) unmarshalEvent(msg *pubsub.Message) (event.Event, error) {
	evIDStr := msg.Attributes[EventIDMessageAttributeKey]
	evID, err := uuid.Parse(evIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event ID '%s': %w", evIDStr, err)
	}
	evName := msg.Attributes[EventNameMessageAttributeKey]
	aggIDStr := msg.Attributes[AggregateIDMessageAttributeKey]

	aggID, err := uuid.Parse(aggIDStr)

	if err != nil {
		return nil, fmt.Errorf("failed to parse aggregate ID '%s': %w", aggIDStr, err)
	}

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

	data, err := a.encoding.Unmarshal(msg.Data, evName)

	if err != nil {
		if a.config.Debug {
			a.logger.Error("Failed to unmarshal event data", zap.Error(err), zap.String("event", eventName), zap.String("data", string(msg.Data)))
		}
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	ev := event.New(
		eventName,
		data,
		event.Aggregate(aggID, aggName, aggVersion),
		event.ID(evID),
		event.Time(evTime),
	)

	return &ev, nil
}

type lifeCycleListener struct {
	adapter *PullAdapter
}

func (l lifeCycleListener) OnStart(ctx context.Context) {
	err := l.adapter.start(ctx)
	if err != nil {
		l.adapter.logger.Error("Failed to start PubSub adapter", zap.Error(err))
		panic(err)
	}
}

func (l lifeCycleListener) OnStop() {
	l.adapter.stop()
}
