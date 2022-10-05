package es

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/AltScore/gothic/pkg/es/event"
	"github.com/stretchr/testify/require"
)

type testSnapshot struct {
	version int
	aValue  string
}

func (t testSnapshot) Apply(event event.Event) error {
	t.aValue = fmt.Sprintf("Applied %s %s", event.Name(), event.ID())
	return nil
}

func (t testSnapshot) SetVersion(version int) {
	t.version = version
}

type testAggregate struct {
	AggregateBase[*testSnapshot]
}

func newTestAgg(events ...event.Event) *testAggregate {
	return &testAggregate{
		AggregateBase: NewAgg[*testSnapshot]("11111", "TestAgg", events, WithSnapshot(&testSnapshot{})),
	}
}

func (a *testAggregate) AddEvent(name string) error {
	return a.Raise(event.For[any](a, name, nil))
}

func Test_new_event_has_correct_version(t *testing.T) {
	agg := newTestAgg()

	err := agg.AddEvent("TestEvent")
	require.NoError(t, err)

	require.Equal(t, 1, agg.Version())
	require.Equal(t, 1, len(agg.GetNewEvents()))
	require.Equal(t, 1, agg.GetNewEvents()[0].Version())
}

func Test_second_new_event_has_correct_version(t *testing.T) {
	agg := newTestAgg()

	require.NoError(t, agg.AddEvent("First event"))
	require.NoError(t, agg.AddEvent("Second Event"))

	require.Equal(t, 2, agg.Version())
	require.Equal(t, 2, len(agg.GetNewEvents()))
	require.Equal(t, 2, agg.GetNewEvents()[1].Version())
}

func Test_replayed_aggregate_has_correct_version(t *testing.T) {
	agg2 := givenAReplayedAggregate(t, 2)

	require.Equal(t, 2, agg2.Version())
	require.Equal(t, 0, len(agg2.GetNewEvents()))
}

func Test_new_event_on_replayed_agg_has_correct_version(t *testing.T) {
	agg := givenAReplayedAggregate(t, 3)

	err := agg.AddEvent("TestEvent")
	require.NoError(t, err)

	require.Equal(t, 4, agg.Version())
	require.Equal(t, 1, len(agg.GetNewEvents()))
	require.Equal(t, 4, agg.GetNewEvents()[0].Version())
}

func givenAReplayedAggregate(t *testing.T, count int) *testAggregate {
	agg := newTestAgg()

	for i := 0; i < count; i++ {
		require.NoError(t, agg.AddEvent("TestEvent "+strconv.Itoa(i)))
	}

	agg2, err := Reify(agg.GetNewEvents(), func(b AggregateBase[*testSnapshot]) *testAggregate {
		return &testAggregate{AggregateBase: b}
	}, WithSnapshot(&testSnapshot{}))

	require.NoError(t, err)

	return agg2
}
