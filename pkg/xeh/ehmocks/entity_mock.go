package ehmocks

import (
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/uuid"
	"github.com/stretchr/testify/mock"
)

type VersionedEntityMock struct {
	mock.Mock
}

var _ eh.Entity = &VersionedEntityMock{}

func (s *VersionedEntityMock) EntityID() uuid.UUID {
	args := s.Called()

	u, ok := args.Get(0).(uuid.UUID)
	if !ok {
		panic("invalid type")
	}
	return u
}

var _ eh.Versionable = &VersionedEntityMock{}

func (s *VersionedEntityMock) AggregateVersion() int {
	args := s.Called()
	return args.Int(0)
}
