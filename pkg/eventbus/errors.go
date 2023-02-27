package eventbus

import (
	"errors"
	"github.com/google/uuid"
)

var (
	// ErrEmptyEventName is returned when an empty event name is provided while subscribing.
	ErrEmptyEventName = errors.New("event name cannot be empty")
	// ErrBusNotRunning is returned when an operation is performed on a localBus that is not running.
	ErrBusNotRunning = errors.New("bus is not running")
	// ErrBusAlreadyRunning is returned when the localBus is already running.
	ErrBusAlreadyRunning = errors.New("bus is already running")
)

type ErrUnhandledEvent struct {
	EventName string
	EventId   uuid.UUID
}

func NewErrUnhandledEvent(eventName string, eventId uuid.UUID) *ErrUnhandledEvent {
	return &ErrUnhandledEvent{EventName: eventName, EventId: eventId}
}

func (e ErrUnhandledEvent) Error() string {
	return "unhandled event: " + e.EventName + " " + e.EventId.String()
}