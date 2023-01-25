package bus

import (
	"context"
	"github.com/modernice/goes/event"
)

// Bus is the pub-sub client for events.
type Bus interface {
	Publisher
	Receiver
}

// A Publisher allows to publish events to subscribers of these events.
type Publisher interface {
	// Publish publishes events. Each event is sent to all subscribers of the event.
	Publish(ctx context.Context, events ...event.Event) error
}

// A Receiver allows to receive events from subscription.
type Receiver interface {
	// Receive subscribes to all events from subscription.
	// To use different handlers for different event type, use a Broker.
	//
	// The provided handler will be called on each received event. If the handler
	// returns an error, the event is not acknowledged and will be received again.
	//
	// When the provided context is canceled, the reception is also canceled.
	//
	// If it cannot start receiving events, an error is returned.
	Receive(ctx context.Context, subscriptionName string, handler EventHandler) error
}
