package eventbus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type EventName = string

type Event interface {
	// ID returns the ID of the event.
	ID() uuid.UUID
	// Name returns the name of the event.
	Name() EventName
}

type EventConsumer func(ctx context.Context, event Event) error
type TypedEventConsumer[T Event] func(ctx context.Context, event T) error

// Callback is a callback function that is called when an event is processed by the listener.
type Callback func(event Event, err error)

type Option func(*EventEnvelope)

// LifeCycleProvider is implemented by buses that allows to add a listener that is called when the bus is started or stopped.
type LifeCycleProvider interface {
	// AddLifecycleListener adds a listener that is called when the bus is started or stopped.
	AddLifecycleListener(listener LifeCycleListener)
}

type LifeCycleListener interface {
	OnStart(ctx context.Context)
	OnStop()
}

type Publisher interface {
	// Publish publishes an event. Methods can return before the Event is handled.
	Publish(event Event, options ...Option) error
}

type Subscriber interface {
	// Subscribe subscribes to an event.
	// The consumer will be called when the event is published.
	// The consumer can return an error, which will be returned by the Publish method.
	Subscribe(eventName EventName, consumer EventConsumer) error
}

// EventBus allows to publish and subscribe to events.
// It is a simple interface that allows to decouple the event publishing from the event handling.
// The handler can return an error, which will be returned by the Publish method.
// The publisher can publish as fire-and-forget, or wait for the handler to confirm the processing (replyHandler).
type EventBus interface {
	// Start starts the event bus to process events.
	Start() error
	// Stop stops the event bus.
	Stop() error

	Publisher
	Subscriber
}

// Subscribe subscribes to an event of a given type. It allows to use a typed handler.
func Subscribe[T Event](subscriber Subscriber, eventName EventName, handler TypedEventConsumer[T]) error {
	untyped := func(ctx context.Context, event Event) error {
		typedEvent, ok := event.(T)
		if !ok {
			var t T
			return fmt.Errorf("event %s/%v type %T is not of type %T", event.Name(), event.ID(), event, t)
		}
		return handler(ctx, typedEvent)
	}

	return subscriber.Subscribe(eventName, untyped)
}
