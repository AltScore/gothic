package es

import (
	"fmt"

	"github.com/AltScore/gothic/pkg/ids"

	"github.com/AltScore/gothic/pkg/es/event"
)

type Snapshot interface {
	Apply(event.Event) error
	SetVersion(int)
}

type Aggregate[SS Snapshot] interface {
	ID() ids.ID
	Type() string
	Version() int
	Snapshot() SS
	Apply(e event.Event) error
	Replay() error
}

type AggregateBase[SS Snapshot] struct {
	type_      string
	id         ids.ID
	version    int
	events     []event.Event
	nextToSave int
	snapshot   SS
}

type Option[SS Snapshot] func(*AggregateBase[SS])

func NewAgg[SS Snapshot](
	id ids.ID,
	type_ string,
	events []event.Event,
	opts ...Option[SS],
) AggregateBase[SS] {
	a := AggregateBase[SS]{
		id:         id,
		type_:      type_,
		events:     events,
		nextToSave: len(events),
	}

	for _, opt := range opts {
		opt(&a)
	}

	return a
}

// Reify recreates an aggregate from a list of events stored to its current state.
func Reify[Agg Aggregate[SS], SS Snapshot](previousEvents []event.Event, factory func(AggregateBase[SS]) Agg, opts ...Option[SS]) (Agg, error) {
	var aggregate Agg
	if len(previousEvents) == 0 {
		return aggregate, fmt.Errorf("no events to rebuild from")
	}

	id, name, _ := previousEvents[0].Aggregate()

	base := NewAgg[SS](id, name, previousEvents, opts...)

	aggregate = factory(base)

	return aggregate, aggregate.Replay()
}

func (a *AggregateBase[SS]) ID() ids.ID {
	return a.id
}

func (a *AggregateBase[SS]) Type() string {
	return a.type_
}

func (a *AggregateBase[SS]) Version() int {
	return a.version
}

func (a *AggregateBase[SS]) Snapshot() SS {
	return a.snapshot
}

func (a *AggregateBase[SS]) Replay() error {
	if a.version > 0 {
		return nil
	}

	for _, e := range a.events {
		if err := a.Apply(e); err != nil {
			return err
		}
	}
	return nil
}

func (a *AggregateBase[SS]) SetId(id ids.ID) {
	a.id = id
}

// Apply process an already existent event to update the current Snapshot
// En error is returned in case the event is incorrect for this Snapshot
func (a *AggregateBase[SS]) Apply(e event.Event) error {
	if err := a.verifyEventCanBeAppliedToThis(e); err != nil {
		return err
	}

	if err := a.snapshot.Apply(e); err != nil {
		return err
	}

	a.version++

	a.snapshot.SetVersion(a.version)

	return nil
}

func (a *AggregateBase[SS]) verifyEventCanBeAppliedToThis(e event.Event) error {
	id, name, version := e.Aggregate()

	if version != a.version+1 {
		return fmt.Errorf("invalid version %d, expected %d", version, a.version)
	}
	if name != a.type_ {
		return fmt.Errorf("invalid entity type %s, expected %s", name, a.type_)
	}

	if a.id != "" && id != a.id {
		return fmt.Errorf("invalid entity id %s, expected %s", id, a.id)
	}
	return nil
}

// Raise process a new event to update the current Snapshot and append it to the past events
func (a *AggregateBase[SS]) Raise(e event.Event) error {
	if err := a.Apply(e); err != nil {
		return err
	}

	a.events = append(a.events, e)
	return nil
}

func (a *AggregateBase[SS]) NewMetadata(eventType string) Metadata {
	return NewMetadata(
		a.type_,
		a.id,
		eventType,
		a.version+1,
	)
}

// Events returns a copy of the events
func (a *AggregateBase[SS]) Events() []event.Event {
	return append([]event.Event{}, a.events...)
}

func (a *AggregateBase[SS]) GetNewEvents() []event.Event {
	return a.events[a.nextToSave:]
}

func (a *AggregateBase[SS]) WithEventsSaved() AggregateBase[SS] {
	return AggregateBase[SS]{
		type_:      a.type_,
		id:         a.id,
		version:    a.version,
		events:     a.events,
		nextToSave: len(a.events),
		snapshot:   a.snapshot,
	}
}

func (a *AggregateBase[SS]) HasEventsToSave() bool {
	return a.nextToSave < len(a.events)
}

func WithSnapshot[SS Snapshot](snapshot SS) Option[SS] {
	return func(a *AggregateBase[SS]) {
		a.snapshot = snapshot
	}
}
