package eventbus

import (
	"github.com/AltScore/gothic/pkg/ids"
)

type EventName = string

type Event interface {
	// ID returns the ID of the event.
	ID() ids.ID
	// Name returns the name of the event.
	Name() EventName
}

type EventHandler func(event Event) error

// Callback is a callback function that is called when an event is processed by the listener.
type Callback func(event Event, err error)

type Option func(*EventEnvelope)

// EventBus allows to publish and subscribe to events.
// It is a simple interface that allows to decouple the event publishing from the event handling.
// The handler can return an error, which will be returned by the Publish method.
// The publisher can publish as fire-and-forget, or wait for the handler to confirm the processing (replyHandler).
type EventBus interface {
	Start() error
	Stop() error
	// Publish publishes an event. Methods can return before the Event is handled.
	Publish(event Event, options ...Option) error

	// Subscribe subscribes to an event.
	// The handler will be called when the event is published.
	// The handler can return an error, which will be returned by the Publish method.
	Subscribe(eventName EventName, handler EventHandler) error
}
