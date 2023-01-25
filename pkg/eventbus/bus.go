package eventbus

import (
	"context"
	"fmt"
	"github.com/AltScore/gothic/pkg/ids"
)

type EventName = string

type Event interface {
	// ID returns the ID of the event.
	ID() ids.ID
	// Name returns the name of the event.
	Name() EventName
}

type EventHandler func(ctx context.Context, event Event) error
type TypedEventHandler[T Event] func(ctx context.Context, event T) error

// Callback is a callback function that is called when an event is processed by the listener.
type Callback func(event Event, err error)

type Option func(*EventEnvelope)

type Publisher interface {
	// Publish publishes an event. Methods can return before the Event is handled.
	Publish(event Event, options ...Option) error
}

type Subscriber interface {
	// Subscribe subscribes to an event.
	// The handler will be called when the event is published.
	// The handler can return an error, which will be returned by the Publish method.
	Subscribe(eventName EventName, handler EventHandler) error
}

// EventBus allows to publish and subscribe to events.
// It is a simple interface that allows to decouple the event publishing from the event handling.
// The handler can return an error, which will be returned by the Publish method.
// The publisher can publish as fire-and-forget, or wait for the handler to confirm the processing (replyHandler).
type EventBus interface {
	Start() error
	Stop() error

	Publisher
	Subscriber
}

// Subscribe subscribes to an event of a given type. It allows to use a typed handler.
func Subscribe[T Event](subscriber Subscriber, eventName EventName, handler TypedEventHandler[T]) error {
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
