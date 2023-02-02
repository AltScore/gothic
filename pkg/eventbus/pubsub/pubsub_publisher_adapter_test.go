package pubsub

import (
	"github.com/AltScore/gothic/pkg/eventbus"
	"github.com/google/uuid"
	"github.com/modernice/goes/event"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Adapter_Given_Generic_event_When_converted_to_event_Then_is_converted(t *testing.T) {
	var originalEvent event.Event = event.New("test", "test").Any()

	var iAmAnEventEvent eventbus.Event = originalEvent

	require.True(t, doCast(iAmAnEventEvent))
}

func Test_Adapter_Given_invalid_event_When_converted_to_event_Then_is_not_converted(t *testing.T) {
	var iAmNotAnEventEvent eventbus.Event = &fakeEvent{}

	require.False(t, doCast(iAmNotAnEventEvent))
}

func doCast(event1 eventbus.Event) bool {
	_, ok := event1.(event.Event)

	return ok
}

type fakeEvent struct {
}

func (f fakeEvent) ID() uuid.UUID {
	panic("implement me")
}

func (f fakeEvent) Name() eventbus.EventName {
	panic("implement me")
}
