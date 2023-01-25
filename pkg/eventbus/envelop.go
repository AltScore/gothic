package eventbus

import (
	"context"
)

// EventEnvelope is a wrapper around the event that is sent to the bus.
// It holds other information that is used by the bus to process the event.
type EventEnvelope struct {
	Event      Event
	Err        error
	Ctx        context.Context
	Callback   Callback
	ShouldWait bool
}

// ProcessOptions processes the options and sets the values on the envelope.
func (e *EventEnvelope) ProcessOptions(options []Option) {
	for _, option := range options {
		option(e)
	}
}

// WithAck sets the callback function that is called when the event is processed.
func WithAck(callback Callback) Option {
	return func(e *EventEnvelope) {
		e.Callback = callback
	}
}

// WithWait sets the flag that indicates that the event should be processed synchronously.
func WithWait() Option {
	return func(e *EventEnvelope) {
		e.ShouldWait = true
	}
}
