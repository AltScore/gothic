package ehmocks

import (
	"context"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/totemcaf/gollections/nils"
)

type ReadRepoMock struct {
	mock.Mock
}

var _ eh.ReadRepo = (*ReadRepoMock)(nil)
var _ eh.WriteRepo = (*ReadRepoMock)(nil)

func (r *ReadRepoMock) InnerRepo(ctx context.Context) eh.ReadRepo {
	args := r.Called(ctx)
	return nils.CastOrNil[eh.ReadRepo](args.Get(0))
}

func (r *ReadRepoMock) Find(ctx context.Context, u uuid.UUID) (eh.Entity, error) {
	args := r.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return nils.CastOrNil[eh.Entity](args.Get(0)), args.Error(1)
}

func (r *ReadRepoMock) FindAll(ctx context.Context) ([]eh.Entity, error) {
	args := r.Called(ctx)
	return nils.CastOrNil[[]eh.Entity](args.Get(0)), args.Error(1)
}

func (r *ReadRepoMock) Close() error {
	args := r.Called()
	return args.Error(0)
}

func (r *ReadRepoMock) Save(ctx context.Context, entity eh.Entity) error {
	args := r.Called(ctx, entity)
	return args.Error(0)
}

func (r *ReadRepoMock) Remove(ctx context.Context, uuid uuid.UUID) error {
	args := r.Called(ctx, uuid)
	return args.Error(0)
}
