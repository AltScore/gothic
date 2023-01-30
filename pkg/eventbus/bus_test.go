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
	// Given a subscriber and type event consumer
	subscriber := &stubSubscriber{}
	called := false

	consumer := func(_ context.Context, event *testTypedEvent) error {
		called = true
		return nil
	}

	// When I subscribe to the event and call the consumer
	err := Subscribe(subscriber, "test", consumer)

	// Then the typed consumer is called
	require.NoError(t, err)
	require.Equal(t, "test", subscriber.eventName)
	require.NotNil(t, subscriber.consumer)
	require.ErrorContains(
		t,
		subscriber.consumer(context.Background(), &anotherTestTypedEvent{}),
		"type *eventbus.anotherTestTypedEvent is not of type *eventbus.testTypedEvent",
	)
	require.False(t, called)
}

func TestSubscribe_a_typed_event_fails_when_called_with_invalid_type(t *testing.T) {
	// Given a subscriber and type event consumer
	subscriber := &stubSubscriber{}
	called := false

	consumer := func(_ context.Context, event *testTypedEvent) error {
		called = true
		return nil
	}

	// When I subscribe to the event and call the consumer
	err := Subscribe(subscriber, "test", consumer)

	// Then the typed consumer is called
	require.NoError(t, err)
	require.Equal(t, "test", subscriber.eventName)
	require.NotNil(t, subscriber.consumer)
	require.NoError(t, subscriber.consumer(context.Background(), &testTypedEvent{}))
	require.True(t, called)
}

type stubSubscriber struct {
	eventName EventName
	consumer  EventConsumer
}

func (s *stubSubscriber) Subscribe(eventName EventName, consumer EventConsumer) error {
	s.eventName = eventName
	s.consumer = consumer
	return nil
}
