package localbus

import (
	"context"
	"fmt"
	"github.com/AltScore/gothic/pkg/eventbus"
	"github.com/AltScore/gothic/pkg/logger"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"strings"
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
	ctx       context.Context
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
	if !b.running.CompareAndSwap(false, true) {
		return eventbus.ErrBusAlreadyRunning
	}

	if b.ctx == nil {
		b.ctx = context.Background()
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
	if !b.running.CompareAndSwap(true, false) {
		return eventbus.ErrBusNotRunning
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
		return eventbus.ErrBusNotRunning
	}

	envelope := &eventbus.EventEnvelope{
		Event: event,
		Ctx:   b.ctx,
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
		return eventbus.ErrEmptyEventName
	}
	b.listeners.addHandler(eventNameTrimmed, handler)
	return nil
}

func (b *localBus) processEvents() {
	b.logger.Info("Starting event bus")
	for {
		select {
		case <-b.ctx.Done():
			b.logger.Info("Stopping event bus due to context done")
			return

		case envelope, ok := <-b.eventCh:
			if !ok {
				b.logger.Info("Stopping event bus due to channel closed")
				return
			}
			b.processEvent(envelope)
		}
	}
}

func (b *localBus) processEvent(envelope *eventbus.EventEnvelope) {
	event := envelope.Event

	listeners := b.listeners.getHandlers(event.Name())

	if len(listeners) == 0 {
		b.logger.Warn(
			"No listeners found for event",
			zap.String("name", event.Name()),
			zap.String("id", event.ID().String()),
		)

		if envelope.Callback != nil {
			// Report unhandled event error to caller
			envelope.Callback(event, eventbus.NewErrUnhandledEvent(event.Name(), event.ID()))
			return
		}

	} else {
		b.logger.Debug(
			"Processing event",
			zap.String("name", event.Name()),
			zap.String("id", event.ID().String()),
		)
	}
	isCallCallbackPending := true

	for _, listener := range listeners {
		if err := b.executeWithRecovery(listener, envelope.Ctx, event); isCallCallbackPending && err != nil {
			// First callback error is the one that will be returned
			// TODO: Should all errors be returned? they can be aggregated into a single error
			if envelope.Callback != nil {
				envelope.Callback(event, envelope.Err)
			}
			isCallCallbackPending = false
		}
	}

	if isCallCallbackPending && envelope.Callback != nil {
		// All callbacks run successfully
		envelope.Callback(event, nil)
	}
}

func (b *localBus) executeWithRecovery(listener eventbus.EventHandler, ctx context.Context, event eventbus.Event) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic while executing event handler: %v", r)
			b.logger.Error("panic while executing event handler", zap.Any("error", r))
		}
	}()

	return listener(ctx, event)
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

func WithContext(ctx context.Context) Option {
	return func(bus *localBus) {
		bus.ctx = ctx
	}
}
