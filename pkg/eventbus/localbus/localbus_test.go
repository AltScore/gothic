package localbus

import (
	"context"
	"fmt"
	"github.com/AltScore/gothic/pkg/eventbus"
	"github.com/AltScore/gothic/pkg/ids"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testEvent struct {
	id   ids.ID
	name string
}

func (t *testEvent) ID() ids.ID {
	return t.id
}

func (t *testEvent) Name() eventbus.EventName {
	return t.name
}

func TestLocalBus_fails_to_publish_if_not_started(t *testing.T) {
	bus := NewLocalBus()
	err := bus.Publish(&testEvent{})

	require.ErrorIs(t, err, eventbus.ErrBusNotRunning)
}

func TestLocalBus_calls_listener(t *testing.T) {
	bus := NewLocalBus(WithBufferSize(1))
	mustStart(t, bus)
	defer mustStop(t, bus)

	called := givenASubscriptionOn("test", bus)

	whenPublishEventOnWithId(t, bus, "test", "ev-1111")

	thenHandlerShouldBeCalled(t, called)
}

func TestLocalBus_does_not_call_listener_for_different_event(t *testing.T) {
	bus := NewLocalBus(WithBufferSize(1))
	mustStart(t, bus)
	defer mustStop(t, bus)

	called := givenASubscriptionOn("test", bus)

	whenPublishEventOnWithId(t, bus, "other", "ev-3658")

	thenHandlerShouldNotBeCalled(t, called)
}

func TestLocalBus_calls_several_listeners(t *testing.T) {
	bus := NewLocalBus(WithBufferSize(1))
	mustStart(t, bus)
	defer mustStop(t, bus)

	called := givenASubscriptionOn("test", bus)
	called2 := givenASubscriptionOn("test", bus)

	whenPublishEventOnWithId(t, bus, "test", "ev-4687")

	// They are called in order
	thenHandlerShouldBeCalled(t, called)
	thenHandlerShouldBeCalled(t, called2)
}

func TestLocalBus_calls_callback_when_event_requires_acknowledge(t *testing.T) {
	// GIVEN a bus with a listener
	bus := NewLocalBus(WithBufferSize(1))
	mustStart(t, bus)
	defer mustStop(t, bus)

	_ = givenASubscriptionOn("test", bus)

	callbackCalled := make(chan string)

	callback := func(event eventbus.Event, err error) {
		callbackCalled <- fmt.Sprintf("%v/%v result: %v", event.Name(), event.ID(), err)
	}

	whenPublishEventOnWithId(t, bus, "test", "ev-1111", callback)

	thenHandlerShouldBeCalled(t, callbackCalled)
}

func TestLocalBus_reports_unhandled_event_error(t *testing.T) {
	// GIVEN a bus without a listener
	bus := NewLocalBus(WithBufferSize(1))
	mustStart(t, bus)
	defer mustStop(t, bus)

	// WHEN publishing an event that is not handled with callback
	callbackCalled := make(chan error, 1)

	callback := func(event eventbus.Event, err error) {
		callbackCalled <- err
	}

	whenPublishEventOnWithId(t, bus, "other", "ev-1111", callback)

	// THEN the callback is called with an error
	thenHandlerShouldBeCalled[error](t, callbackCalled, eventbus.NewErrUnhandledEvent("other", "ev-1111"))
}

func givenASubscriptionOn(name string, bus eventbus.EventBus) chan ids.ID {
	called := make(chan ids.ID, 1) // Need a buffered channel to avoid blocking
	_ = bus.Subscribe(name, func(_ context.Context, event eventbus.Event) error {
		fmt.Printf("Received event: %v/%v\n", event.Name(), event.ID())
		called <- event.ID()
		return nil
	})
	return called
}

func whenPublishEventOnWithId(t *testing.T, bus eventbus.EventBus, name string, id ids.ID, callbacks ...eventbus.Callback) {
	options := make([]eventbus.Option, 0)
	for _, callback := range callbacks {
		options = append(options, eventbus.WithAck(callback))
	}
	event := &testEvent{name: name, id: id}

	fmt.Printf("Publishing event %v/%v\n", event.Name(), event.ID())

	err := bus.Publish(event, options...)
	require.NoError(t, err)
}

func thenHandlerShouldBeCalled[T any](t *testing.T, called chan T, expected ...T) {
	// wait for channel called to be processed with timeout
	fmt.Println("Waiting for channel called to be processed")
	select {
	case result := <-called:
		// ok
		if len(expected) > 0 {
			require.Equal(t, expected[0], result)
		}
		fmt.Printf("Handler was called with %v\n", result)
	case <-time.After(10 * time.Millisecond):
		require.Fail(t, "Listener not called")
	}
}

func thenHandlerShouldNotBeCalled[T any](t *testing.T, called chan T) {
	// wait for channel called to be processed with timeout
	fmt.Println("Waiting for channel called to be processed")
	select {
	case <-called:
		require.Fail(t, "Listener was called but not expected")

	case <-time.After(10 * time.Millisecond):
		// ok
		fmt.Println("Handler was not called -> Ok")
	}
}

func mustStart(t *testing.T, bus eventbus.EventBus) {
	err := bus.Start()
	require.NoError(t, err)
}

func mustStop(t *testing.T, bus eventbus.EventBus) {
	err := bus.Stop()
	require.NoError(t, err)
}
