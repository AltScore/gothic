package broker

import (
	"context"
	"fmt"
	"github.com/modernice/goes/event"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func Test_Broker_calls_subscribed_handler_for_new_message(t *testing.T) {
	// GIVEN new broker
	br := New(zap.NewNop())

	// AND a handler was subscribed to the event
	called := false
	handler := func(_ context.Context, event event.Event) error {
		called = true
		return nil
	}

	_, _ = br.Subscribe(handler, "event1")

	// WHEN a message is handled
	ev := event.New("event1", "a sample value")

	err := br.Handle(context.Background(), ev.Any())

	// THEN the handler is called
	require.True(t, called)

	// AND no error is returned
	require.NoError(t, err)
}

func Test_Broker_calls_subscribed_handler_for_new_message_with_multiple_subscriptions(t *testing.T) {
	// GIVEN new broker
	br := New(zap.NewNop())

	// AND a handler was subscribed to the event
	called := ""
	handler := func(_ context.Context, ev event.Event) error {

		tev := event.Cast[string, any](ev)

		called = tev.Data()
		return nil
	}

	_, _ = br.Subscribe(handler, "event1", "event2")

	// WHEN a message is handled for event1
	ev := event.New("event1", "a sample value")

	err := br.Handle(context.Background(), ev.Any())

	// THEN the handler is called with correct value
	require.Equal(t, "a sample value", called)

	// AND no error is returned
	require.NoError(t, err)

	// WHEN a message is handled for event2
	ev2 := event.New("event2", "another sample value")

	err2 := br.Handle(context.Background(), ev2.Any())

	// THEN the handler is called
	require.Equal(t, "another sample value", called)

	// AND no error is returned
	require.NoError(t, err2)
}

func Test_Broker_does_not_call_subscribed_handler_for_different_event(t *testing.T) {
	// GIVEN new broker
	br := New(zap.NewNop())

	// AND a handler was subscribed to the event
	called := false
	handler := func(_ context.Context, event event.Event) error {
		called = true
		return nil
	}

	_, _ = br.Subscribe(handler, "event1")

	// WHEN a message is handled
	ev := event.New("event2", "a sample value")

	err := br.Handle(context.Background(), ev.Any())

	// THEN the handler is not called
	require.False(t, called)

	// AND no error is returned
	require.NoError(t, err)
}

func Test_Broker_calls_all_subscribed_handlers_for_new_message(t *testing.T) {
	// GIVEN new broker
	br := New(zap.NewNop())

	// AND a handler was subscribed to the event
	called1 := 0
	handler1 := func(_ context.Context, event event.Event) error {
		called1++
		return nil
	}

	_, _ = br.Subscribe(handler1, "event1")

	// AND another handler was subscribed to the event
	called2 := 0
	handler2 := func(_ context.Context, event event.Event) error {
		called2++
		return nil
	}

	_, _ = br.Subscribe(handler2, "event1")

	// WHEN a message is handled
	ev := event.New("event1", "a sample value")

	err := br.Handle(context.Background(), ev.Any())

	// THEN all the handlers are called once
	require.Equal(t, 1, called1)
	require.Equal(t, 1, called2)

	// AND no error is returned
	require.NoError(t, err)
}

func Test_Broker_calls_all_subscribed_handlers_for_new_message_2(t *testing.T) {
	// GIVEN new broker
	br := New(zap.NewNop())

	// AND a handler was subscribed to the event
	called := 0
	handler := func(_ context.Context, event event.Event) error {
		called++
		return nil
	}

	_, _ = br.Subscribe(handler, "event1")
	_, _ = br.Subscribe(handler, "event1")

	// WHEN a message is handled
	ev := event.New("event1", "a sample value")

	err := br.Handle(context.Background(), ev.Any())

	// THEN the handler is called twice
	require.Equal(t, 2, called)

	// AND no error is returned
	require.NoError(t, err)
}

func Test_Broker_returns_error_for_first_failing_handler(t *testing.T) {
	errTest := fmt.Errorf("test error")

	// GIVEN new broker
	br := New(zap.NewNop())

	// AND a handler was subscribed to the event
	called := 0
	handler := func(_ context.Context, event event.Event) error {
		called++
		return errTest
	}

	_, _ = br.Subscribe(handler, "event1")

	// AND another handler was subscribed to the event
	called2 := 0
	handler2 := func(_ context.Context, event event.Event) error {
		called2++
		return nil
	}

	_, _ = br.Subscribe(handler2, "event1")

	// WHEN a message is handled
	ev := event.New("event1", "a sample value")

	err := br.Handle(context.Background(), ev.Any())

	// THEN the handler is called once
	require.Equal(t, 1, called)

	// AND the second handler is called once
	require.Equal(t, 1, called2)

	// AND the error is returned
	require.ErrorContains(t, err, errTest.Error())
}

func Test_Broker_returns_error_for_first_failing_handler_2(t *testing.T) {
	errTest := fmt.Errorf("test error")

	// GIVEN new broker
	br := New(zap.NewNop())

	// AND a handler was subscribed to the event
	called := 0
	handler := func(_ context.Context, event event.Event) error {
		called++
		return nil
	}

	_, _ = br.Subscribe(handler, "event1")

	// AND another handler was subscribed to the event
	called2 := 0
	handler2 := func(_ context.Context, event event.Event) error {
		called2++
		return errTest
	}

	_, _ = br.Subscribe(handler2, "event1")

	// WHEN a message is handled
	ev := event.New("event1", "a sample value")

	err := br.Handle(context.Background(), ev.Any())

	// THEN the handler is called once
	require.Equal(t, 1, called)

	// AND the second handler is called once
	require.Equal(t, 1, called2)

	// AND the error is returned
	require.ErrorContains(t, err, errTest.Error())
}

func Test_Broker_waits_all_subscribed_handlers_to_finish(t *testing.T) {
	// GIVEN new broker
	br := New(zap.NewNop())

	// AND a handler was subscribed to the event
	called1 := 0
	handler1 := func(_ context.Context, event event.Event) error {
		called1++
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	_, _ = br.Subscribe(handler1, "event1")

	// AND another handler was subscribed to the event
	called2 := 0
	handler2 := func(_ context.Context, event event.Event) error {
		called2++
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	_, _ = br.Subscribe(handler2, "event1")

	// WHEN a message is handled
	ev := event.New("event1", "a sample value")

	err := br.Handle(context.Background(), ev.Any())

	// THEN all the handlers are called once
	require.Equal(t, 1, called1)
	require.Equal(t, 1, called2)

	// AND no error is returned
	require.NoError(t, err)
}

func Test_Broker_waits_all_subscribed_handlers_to_finish_2(t *testing.T) {
	// GIVEN new broker
	br := New(zap.NewNop())

	// AND a handler was subscribed to the event
	called := 0
	handler := func(_ context.Context, event event.Event) error {
		called++
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	_, _ = br.Subscribe(handler, "event1")
	_, _ = br.Subscribe(handler, "event1")

	// WHEN a message is handled
	ev := event.New("event1", "a sample value")

	err := br.Handle(context.Background(), ev.Any())

	// THEN the handler is called twice
	require.Equal(t, 2, called)

	// AND no error is returned
	require.NoError(t, err)
}

func Compare2Arrays(t *testing.T, expected, actual []string) {
	require.Equal(t, len(expected), len(actual))
	for i := range expected {
		require.Equal(t, expected[i], actual[i])
	}
}
