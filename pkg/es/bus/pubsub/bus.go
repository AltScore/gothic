package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AltScore/gothic/pkg/es/bus"
	l "github.com/AltScore/gothic/pkg/logger"
	"github.com/google/uuid"
	"github.com/modernice/goes/codec"
	"github.com/modernice/goes/event"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

const (
	AggregateIDMessageAttributeKey      = "aggID"
	AggregateNameMessageAttributeKey    = "aggName"
	AggregateVersionMessageAttributeKey = "aggVer"
	EventIDMessageAttributeKey          = "id"
	EventNameMessageAttributeKey        = "name"
	EventTimeMessageAttributeKey        = "time"
)

type Config struct {
	ProjectID string `mapstructure:"project_id"`
	TopicName string `mapstructure:"topic_name"`
}

type handlerEntry struct {
	eventName string
	handler   bus.EventHandler
	enabled   bool // if false, no events are sent to it
}

type Bus struct {
	logger l.Logger
	client *pubsub.Client
	topic  *pubsub.Topic

	lock sync.RWMutex // control concurrent access to subscriptions and publishers

	publishers    map[string]bus.Publisher
	subscriptions map[string]*pubsub.Subscription
	registry      *codec.Registry
}

func NewBus(logger l.Logger, config Config, registry *codec.Registry) *Bus {
	if config.TopicName == "" {
		panic("PubSub topic name is required")
	}
	if config.ProjectID == "" {
		panic("PubSub project ID is required")
	}

	client, err := pubsub.NewClient(context.Background(), config.ProjectID)
	if err != nil {
		logger.Fatal("Cannot connect to PubSub", zap.Error(err))
		panic(err)
	}

	logger.Info("Connected to PubSub", zap.String("project_id", config.ProjectID), zap.String("topic_name", config.TopicName))

	topic := client.Topic(config.TopicName)
	topic.EnableMessageOrdering = true // This is required for the OrderingKey to work. It is critical for Aggregate Event Sourcing

	return &Bus{
		logger:        logger,
		client:        client,
		topic:         topic,
		publishers:    make(map[string]bus.Publisher),
		subscriptions: make(map[string]*pubsub.Subscription, 0),
		registry:      registry,
	}
}

func (b *Bus) GetPublisher(topicName string) bus.Publisher {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if p, ok := b.publishers[topicName]; ok {
		return p
	}

	p := &publisher{b.logger, b.client.Topic(topicName)}
	b.publishers[topicName] = p

	return p
}

type publisher struct {
	logger l.Logger
	topic  *pubsub.Topic
}

// Publish send the given events to the configured PubSub topic
// Each message is sent in order, if an error is produced, it stops sending and returns the error
func (p *publisher) Publish(ctx context.Context, events ...event.Event) error {
	for _, e := range events {
		data, err := json.Marshal(e)
		if err != nil {
			return err
		}

		aggID, _, _ := e.Aggregate()

		result := p.topic.Publish(ctx, &pubsub.Message{
			Data:        data,
			OrderingKey: aggID.String(),
		})

		_, err = result.Get(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bus) Receive(ctx context.Context, subscriptionName string, handler bus.EventHandler) error {
	s := b.findOrCreateSubscription(subscriptionName)

	if err := s.Receive(ctx, b.makeMsgReceiver(handler)); err != nil {
		b.logger.Warn("error processing events", zap.String("subscription", s.String()), zap.Error(err))
	}
	return nil
}

func (b *Bus) findOrCreateSubscription(name string) *pubsub.Subscription {
	// This method is not going to be called often, so, using
	// a write-lock simplifies escalation.
	b.lock.Lock()
	defer b.lock.Unlock()

	if s, ok := b.subscriptions[name]; ok {
		return s
	}

	s := b.client.Subscription(name)
	b.subscriptions[name] = s
	return s
}

func (b *Bus) makeMsgReceiver(handler bus.EventHandler) func(ctx context.Context, message *pubsub.Message) {
	return func(ctx context.Context, message *pubsub.Message) {
		shouldAck := true
		defer func() {
			if shouldAck {
				message.Ack()
			} else {
				message.Nack()
			}
		}()

		ev, err := b.getEvent(message)
		if err != nil {
			b.logger.Warn("error unmarshalling event", zap.String("orderingKey", message.OrderingKey), zap.String("id", message.ID))
			// cannot process it
			return
		}

		if err := handler(ctx, ev); err != nil {
			b.logger.Warn(
				"error processing event",
				zap.String("orderingKey", message.OrderingKey),
				zap.String("id", message.ID),
				zap.String("type", ev.Name()),
				zap.Error(err),
			)
			// cannot process it
			return
		}
	}
}

func (b *Bus) getEvent(message *pubsub.Message) (event.Event, error) {

	eventType := message.Attributes[EventNameMessageAttributeKey]

	if eventType == "" {
		return nil, fmt.Errorf("event type not found in message attributes (%v)", message.ID)
	}

	data, err := b.registry.New(eventType)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(message.Data, data); err != nil {
		return nil, err
	}

	aggregateOpt, err := b.getAggregateOption(message)
	if err != nil {
		return nil, err
	}

	idOpt, err := b.getIDOption(message)
	if err != nil {
		return nil, err
	}

	timeOpt, err := b.getTimeOption(message)
	if err != nil {
		return nil, err
	}

	ev := event.New(
		eventType,
		data,
		aggregateOpt,
		idOpt,
		timeOpt,
	)

	return ev, nil
}

func (b *Bus) getIDOption(message *pubsub.Message) (event.Option, error) {
	idStr := message.Attributes[EventIDMessageAttributeKey]
	id, err := uuid.Parse(idStr)

	if err != nil {
		return DoNoting(), fmt.Errorf("invalid event ID '%s': %s", idStr, err)
	}
	return event.ID(id), nil
}

func (b *Bus) getAggregateOption(message *pubsub.Message) (event.Option, error) {
	attributes := message.Attributes

	name := attributes[AggregateNameMessageAttributeKey]

	if name == "" {
		return DoNoting(), nil
	}

	aggIDStr := attributes[AggregateIDMessageAttributeKey]

	aggID, err := uuid.Parse(aggIDStr)

	if err != nil {
		return DoNoting(), fmt.Errorf("invalid aggregate ID '%s': %s", aggIDStr, err)
	}

	aggVersionStr := attributes[AggregateVersionMessageAttributeKey]
	aggVersion, err := strconv.Atoi(aggVersionStr)

	if err != nil {
		return DoNoting(), fmt.Errorf("invalid version '%s': %v", aggVersionStr, err)
	}

	aggregate := event.Aggregate(
		aggID,
		name,
		aggVersion,
	)
	return aggregate, nil
}

func (b *Bus) getTimeOption(message *pubsub.Message) (event.Option, error) {
	timeStr := message.Attributes[EventTimeMessageAttributeKey]

	if timeStr == "" {
		return DoNoting(), nil
	}

	tm, err := time.Parse(time.RFC3339, timeStr)

	if err != nil {
		b.logger.Warn("error parsing event time", zap.String("time", timeStr), zap.Error(err), zap.String("id", message.ID))
		tm = message.PublishTime
	}

	return event.Time(tm), nil
}

func DoNoting() event.Option {
	return func(evt *event.Evt[any]) {
		// do nothing
	}
}
