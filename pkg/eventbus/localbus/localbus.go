package localbus

import (
	"errors"
	"github.com/AltScore/gothic/pkg/eventbus"
	"github.com/AltScore/gothic/pkg/logger"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"strings"
)

var (
	// ErrEmptyEventName is returned when an empty event name is provided while subscribing.
	ErrEmptyEventName = errors.New("event name cannot be empty")
	// ErrBusNotRunning is returned when an operation is performed on a localBus that is not running.
	ErrBusNotRunning = errors.New("bus is not running")
	// ErrBusAlreadyRunning is returned when the localBus is already running.
	ErrBusAlreadyRunning = errors.New("bus is already running")
)

// Option is a function that configures a localBus
type Option func(bus *localBus)

// localBus is an in-memory event bus implementation of EventBus interface.
type localBus struct {
	logger    logger.Logger
	listeners handlersMap
	eventCh   chan *eventbus.EventEnvelope
	running   atomic.Bool
	size      int
}

// NewLocalBus creates a new localBus instance and configures it with the given options.
func NewLocalBus(options ...Option) eventbus.EventBus {
	bus := &localBus{
		listeners: handlersMap{},
	}

	bus.applyOptions(options)

	return bus
}

func (b *localBus) applyOptions(options []Option) {
	for _, option := range options {
		option(b)
	}
}

// Start starts the localBus event processing loop.
// It is safe to call Start multiple times, but only the first call will start the event processing loop.
func (b *localBus) Start() error {
	if !b.running.CAS(false, true) {
		return ErrBusAlreadyRunning
	}

	if b.logger == nil {
		b.logger = logger.New()
	}

	b.eventCh = make(chan *eventbus.EventEnvelope, b.size)
	go b.processEvents()

	return nil
}

// Stop stops the localBus event processing loop.
// It is safe to call Stop multiple times, but only the first call will stop the event processing loop.
// After stopping the event processing loop, the localBus can be started again.
func (b *localBus) Stop() error {
	if !b.running.CAS(true, false) {
		return ErrBusNotRunning
	}

	close(b.eventCh)
	return nil
}

// Publish publishes an event to the localBus.
// All the registered handlers for the event will be called in a separate goroutine.
// Publish is non-blocking and returns immediately.
// If any of the handlers returns an error, the error will be returned by Publish to the configured callback.
func (b *localBus) Publish(event eventbus.Event, options ...eventbus.Option) error {
	if !b.running.Load() {
		return ErrBusNotRunning
	}

	envelope := &eventbus.EventEnvelope{
		Event: event,
	}

	envelope.ProcessOptions(options)

	b.eventCh <- envelope

	if envelope.ShouldWait {
		// TODO
	}

	return nil
}

// Subscribe subscribes a handler to an event.
// The handler will be called when an event with the given name is published.
// The handler will be called in a separate goroutine.
// If the handler returns an error, the error will be returned by Publish to the configured callback.
// All the handlers for an event will be called in the order they were registered.
func (b *localBus) Subscribe(eventName eventbus.EventName, handler eventbus.EventHandler) error {
	eventNameTrimmed := strings.TrimSpace(eventName)
	if eventNameTrimmed == "" {
		return ErrEmptyEventName
	}
	b.listeners.addHandler(eventNameTrimmed, handler)
	return nil
}

func (b *localBus) processEvents() {
	b.logger.Info("Starting event bus")
	for envelope := range b.eventCh {
		b.logger.Debug(
			"Processing event",
			zap.String("name", envelope.Event.Name()),
			zap.String("id", envelope.Event.ID().String()),
		)
		listeners := b.listeners.getHandlers(envelope.Event.Name())
		isCallCallbackPending := true
		for _, listener := range listeners {
			if err := listener(envelope.Event); isCallCallbackPending && err != nil {
				// First callback error is the one that will be returned
				// TODO: Should all errors be returned? they can be aggregated into a single error
				if envelope.Callback != nil {
					envelope.Callback(envelope.Event, envelope.Err)
				}
				isCallCallbackPending = false
			}

		}

		if isCallCallbackPending && envelope.Callback != nil {
			// All callbacks run successfully
			envelope.Callback(envelope.Event, nil)
		}
	}
}

// WithBufferSize sets the size of the event buffer.
// If the buffer is full, Publish will block until there is space in the buffer.
func WithBufferSize(size int) Option {
	return func(bus *localBus) {
		bus.size = size
	}
}

// WithLogger sets the logger to use by the localBus.
func WithLogger(logger logger.Logger) Option {
	return func(bus *localBus) {
		bus.logger = logger
	}
}
