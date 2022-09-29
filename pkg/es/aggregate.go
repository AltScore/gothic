package es

import (
	"fmt"

	"github.com/AltScore/gothic/pkg/es/event"
)

type Snapshot interface {
	Apply(event.Event) error
	SetVersion(int)
}

type AggregateBase[SS Snapshot] struct {
	type_      string
	id         string
	version    int
	events     []event.Event
	nextToSave int
	snapshot   SS
}

type Option[SS Snapshot] func(*AggregateBase[SS])

func NewAgg[SS Snapshot](
	id string,
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

func (a *AggregateBase[SS]) ID() string {
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
	for _, e := range a.events {
		if err := a.snapshot.Apply(e); err != nil {
			return err
		}
	}
	return nil
}

func (a *AggregateBase[SS]) SetId(id string) {
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

func (a *AggregateBase[SS]) GetNewEvents() []event.Event {
	return a.events[a.nextToSave:]
}

func (a *AggregateBase[SS]) MarkEventsAsSaved() {
	a.nextToSave = len(a.events)
}

func (a *AggregateBase[SS]) HasEventsToSave() bool {
	return a.nextToSave < len(a.events)
}

func WithSnapshot[SS Snapshot](snapshot SS) Option[SS] {
	return func(a *AggregateBase[SS]) {
		a.snapshot = snapshot
	}
}
