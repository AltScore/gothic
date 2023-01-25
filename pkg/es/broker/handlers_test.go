package broker

import (
	"context"
	"errors"
	"github.com/modernice/goes/event"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_new_Handlers_has_no_events(t *testing.T) {
	newHandlers := Handlers{}

	require.Len(t, newHandlers, 0)
}

func Test_Handlers_calls_added_handler(t *testing.T) {
	newHandlers := Handlers{}
	var eventReceived event.Event

	newHandlers.RegisterEventHandler("event1", func(ctx context.Context, event event.Event) error {
		eventReceived = event
		return nil
	})

	// WHEN we call the handler
	eventSent := event.New("event1", "data1").Any()
	err := newHandlers.HandleEvent(context.Background(), eventSent)

	// THEN no error is returned
	require.NoError(t, err)

	// THEN the handler should be called
	require.Equal(t, eventSent, eventReceived)
}

func Test_Handlers_does_not_call_added_handler_for_different_event_name(t *testing.T) {
	newHandlers := Handlers{}
	var eventReceived event.Event

	newHandlers.RegisterEventHandler("event1", func(ctx context.Context, event event.Event) error {
		eventReceived = event
		return nil
	})

	// WHEN we call the handler
	eventSent := event.New("other_event", "data1").Any()
	err := newHandlers.HandleEvent(context.Background(), eventSent)

	// THEN no error is returned
	require.NoError(t, err)

	// THEN the handler should be called
	require.Nil(t, eventReceived)
}

func Test_Handlers_returns_error_of_first_handler_that_fails(t *testing.T) {
	newHandlers := Handlers{}

	newHandlers.RegisterEventHandler("event", func(ctx context.Context, event event.Event) error {
		return errors.New("error 1")
	})

	newHandlers.RegisterEventHandler("event", func(ctx context.Context, event event.Event) error {
		return errors.New("error 2")
	})

	// WHEN we call the handler
	eventSent := event.New("event", "data1").Any()
	err := newHandlers.HandleEvent(context.Background(), eventSent)

	// THEN returned error is the error of the first handler
	require.ErrorContains(t, err, "error 1")
}

func Test_Handlers_does_not_call_handler_2_if_1_fails(t *testing.T) {
	newHandlers := Handlers{}

	called := false

	newHandlers.RegisterEventHandler("event", func(ctx context.Context, event event.Event) error {
		return errors.New("error 1")
	})

	newHandlers.RegisterEventHandler("event", func(ctx context.Context, event event.Event) error {
		called = true
		return nil
	})

	// WHEN we call the handler
	eventSent := event.New("event", "data1").Any()
	_ = newHandlers.HandleEvent(context.Background(), eventSent)

	// THEN handler 2 is not called
	require.False(t, called)
}
