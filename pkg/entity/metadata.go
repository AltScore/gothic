package entity

import (
	"time"

	"github.com/AltScore/gothic/pkg/ids"
)

// Metadata is the metadata for any entity.
type Metadata struct {
	ID        ids.ID    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func New() Metadata {
	return Metadata{
		ID: ids.New(),
	}
}

func NewAt(now time.Time) Metadata {
	return Metadata{
		ID:        ids.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GetID returns the ID of the entity. Implements Identifiable interfaces.
func (e Metadata) GetID() ids.ID {
	return e.ID
}

// Clone returns a clone of the entity with a new ID and CreatedAt if necessary. Updates UpdatedAt.
func (e Metadata) Clone(now time.Time) Metadata {
	return Metadata{
		ID:        e.ID.SelfOrNew(),
		CreatedAt: e.createdAtOrNow(now),
		UpdatedAt: now,
	}
}

func (e Metadata) createdAtOrNow(now time.Time) time.Time {
	var created = e.CreatedAt
	if !created.IsZero() {
		return created
	}
	return now
}
