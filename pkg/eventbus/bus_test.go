package eventbus

import (
	"context"
	"github.com/AltScore/gothic/pkg/ids"
	"github.com/stretchr/testify/require"
	"testing"
)

type testTypedEvent struct {
	id   ids.ID
	name string
}

func (t *testTypedEvent) ID() ids.ID {
	return t.id
}

func (t *testTypedEvent) Name() EventName {
	return t.name
}

type anotherTestTypedEvent struct {
	testTypedEvent
}

func TestSubscribe_a_typed_event(t *testing.T) {
	// Given a subscriber and type event handler
	subscriber := &stubSubscriber{}
	called := false

	handler := func(_ context.Context, event *testTypedEvent) error {
		called = true
		return nil
	}

	// When I subscribe to the event and call the handler
	err := Subscribe(subscriber, "test", handler)

	// Then the typed handler is called
	require.NoError(t, err)
	require.Equal(t, "test", subscriber.eventName)
	require.NotNil(t, subscriber.handler)
	require.ErrorContains(
		t,
		subscriber.handler(context.Background(), &anotherTestTypedEvent{}),
		"type *eventbus.anotherTestTypedEvent is not of type *eventbus.testTypedEvent",
	)
	require.False(t, called)
}

func TestSubscribe_a_typed_event_fails_when_called_with_invalid_type(t *testing.T) {
	// Given a subscriber and type event handler
	subscriber := &stubSubscriber{}
	called := false

	handler := func(_ context.Context, event *testTypedEvent) error {
		called = true
		return nil
	}

	// When I subscribe to the event and call the handler
	err := Subscribe(subscriber, "test", handler)

	// Then the typed handler is called
	require.NoError(t, err)
	require.Equal(t, "test", subscriber.eventName)
	require.NotNil(t, subscriber.handler)
	require.NoError(t, subscriber.handler(context.Background(), &testTypedEvent{}))
	require.True(t, called)
}

type stubSubscriber struct {
	eventName EventName
	handler   EventHandler
}

func (s *stubSubscriber) Subscribe(eventName EventName, handler EventHandler) error {
	s.eventName = eventName
	s.handler = handler
	return nil
}
