package es

import (
	"github.com/AltScore/gothic/pkg/es/event"
)

// EventSource is an interface implemented by classes that can provide events.
type EventSource interface {
	Events() []event.Event
}
