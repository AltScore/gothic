package es

import (
	"fmt"
)

type SnapshotWithVersion interface {
	SetVersion(int)
}

type Versioned interface {
	SetVersion(int)
}

type AggregateBase[ID EntityID[ID], Snapshot Versioned] struct {
	entityType string
	entityID   ID
	version    int
	events     []Event[ID, Snapshot]
	nextToSave int
	snapshot   Snapshot
}

func NewAgg[ID EntityID[ID], Snapshot Versioned](
	entityType string,
	events []Event[ID, Snapshot],
) AggregateBase[ID, Snapshot] {
	return AggregateBase[ID, Snapshot]{
		entityType: entityType,
		events:     events,
		nextToSave: len(events),
	}
}

func (a *AggregateBase[ID, Snapshot]) EntityType() string {
	return a.entityType
}

func (a *AggregateBase[ID, Snapshot]) Version() int {
	return a.version
}

func (a *AggregateBase[ID, Snapshot]) Snapshot() *Snapshot {
	return &a.snapshot
}

func (a *AggregateBase[ID, Snapshot]) Replay() error {
	if len(a.events) == 0 {
		a.entityID = a.entityID.New()
	} else {
		a.entityID = a.events[0].EntityID()
	}

	for _, e := range a.events {
		if err := e.Apply(&a.snapshot); err != nil {
			return err
		}
	}
	return nil
}

func (a *AggregateBase[ID, Snapshot]) SetId(id ID) {
	a.entityID = id
}

// Apply process an already existent event to update the current Snapshot
// En error is returned in case the event is incorrect for this Snapshot
func (a *AggregateBase[ID, Snapshot]) Apply(e Event[ID, Snapshot]) error {
	if err := a.verifyEventCanBeAppliedToThis(e); err != nil {
		return err
	}

	if err := e.Apply(&a.snapshot); err != nil {
		return err
	}

	a.version = e.Version()

	a.snapshot.SetVersion(a.version)

	return nil
}

func (a *AggregateBase[ID, Snapshot]) verifyEventCanBeAppliedToThis(e Event[ID, Snapshot]) error {
	if e.Version() != a.version+1 {
		return fmt.Errorf("invalid version %d, expected %d", e.Version(), a.version)
	}
	if a.entityType != e.EntityType() {
		return fmt.Errorf("invalid entity type %s, expected %s", e.EntityType(), a.entityType)
	}

	if !a.entityID.Empty() && !a.entityID.Eq(e.EntityID()) {
		return fmt.Errorf("invalid entity id %s, expected %s", e.EntityID(), a.entityID)
	}
	return nil
}

// Raise process a new event to update the current Snapshot and append it to the past events
func (a *AggregateBase[ID, Snapshot]) Raise(e Event[ID, Snapshot]) error {
	if err := a.Apply(e); err != nil {
		return err
	}

	a.events = append(a.events, e)
	return nil
}

func (a *AggregateBase[ID, Snapshot]) NewMetadata(eventType string) Metadata[ID] {
	return NewMetadata[ID](
		a.entityType,
		a.entityID,
		eventType,
		a.version+1,
	)
}

func (a *AggregateBase[ID, Snapshot]) GetNewEvents() []Event[ID, Snapshot] {
	return a.events[a.nextToSave:]
}
