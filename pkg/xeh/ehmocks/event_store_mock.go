package ehmocks

import (
	"context"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/totemcaf/gollections/nils"
)

type EventStoreMock struct {
	mock.Mock
}

var _ eh.EventStore = (*EventStoreMock)(nil)

func (e *EventStoreMock) LoadFrom(ctx context.Context, id uuid.UUID, version int) ([]eh.Event, error) {
	args := e.Called(ctx, id, version)
	return nils.CastOrNil[[]eh.Event](args.Get(0)), args.Error(1)
}

func (e *EventStoreMock) Save(ctx context.Context, events []eh.Event, originalVersion int) error {
	args := e.Called(ctx, events, originalVersion)
	return args.Error(0)
}

func (e *EventStoreMock) Load(ctx context.Context, u uuid.UUID) ([]eh.Event, error) {
	args := e.Called(ctx, u)
	return nils.CastOrNil[[]eh.Event](args.Get(0)), args.Error(1)
}

func (e *EventStoreMock) Close() error {
	args := e.Called()
	return args.Error(0)
}

type EntityFake struct {
	ID uuid.UUID
}

var _ eh.Entity = (*EntityFake)(nil)

func (e *EntityFake) EntityID() uuid.UUID {
	return e.ID
}
