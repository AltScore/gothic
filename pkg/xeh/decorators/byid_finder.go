package decorators

import (
	"context"
	"fmt"
	"github.com/AltScore/gothic/v2/pkg/ids"
	eh "github.com/looplab/eventhorizon"
)

// ByIDFinder is a decorator for ReadRepo that adds a FindById method.
type ByIDFinder[Entity eh.Entity] struct {
	eh.ReadRepo
	empty Entity
}

// NewByIdFinder decorates a read repo with a FindById method.
func NewByIdFinder[Entity eh.Entity](readRepo eh.ReadRepo) *ByIDFinder[Entity] {
	return &ByIDFinder[Entity]{ReadRepo: readRepo}
}

func (r *ByIDFinder[Entity]) FindById(ctx context.Context, id ids.Id) (Entity, error) {
	ent, err := r.ReadRepo.Find(ctx, id)

	if err != nil {
		return r.empty, err
	}

	typedEntity, ok := ent.(Entity)

	if !ok {
		return r.empty, fmt.Errorf("entity %v is not of type %T", id, typedEntity)
	}

	return typedEntity, nil
}
