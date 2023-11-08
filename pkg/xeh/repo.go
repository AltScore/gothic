package xeh

import (
	"context"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/uuid"
)

type ReadRepo[Entity eh.Entity] interface {
	// InnerRepo returns the inner read repository, if there is one.
	// Useful for iterating a wrapped set of repositories to get a specific one.
	InnerRepo(context.Context) eh.ReadRepo

	// Find returns an entity for an Id.
	Find(context.Context, uuid.UUID) (Entity, error)

	// FindAll returns all entities in the repository.
	FindAll(context.Context) ([]Entity, error)

	// Close closes the ReadRepo.
	Close() error
}
