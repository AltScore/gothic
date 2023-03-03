package entity

import (
	"context"
	"github.com/AltScore/gothic/pkg/xcontext"
	"time"

	"github.com/AltScore/gothic/pkg/ids"
)

// Metadata is the metadata for any entity.
type Metadata struct {
	ID        ids.ID    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Version   int       `json:"version"`
	TenantID  string    `json:"tenant"`
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

// NewIn creates a new entity with the given context.
func NewIn(ctx context.Context) Metadata {
	tenantID := xcontext.TenantOrDefault(ctx)

	return Metadata{
		ID:       ids.New(),
		TenantID: tenantID,
	}
}

// NewInAt creates a new entity with the given context and time.
func NewInAt(ctx context.Context, now time.Time) Metadata {
	tenantID := xcontext.TenantOrDefault(ctx)

	return Metadata{
		ID:        ids.New(),
		CreatedAt: now,
		UpdatedAt: now,
		TenantID:  tenantID,
	}
}

// GetID returns the ID of the entity. Implements Identifiable interfaces.
func (e Metadata) GetID() ids.ID {
	return e.ID
}

func (e Metadata) Tenant() string {
	return e.TenantID
}

// Clone returns a clone of the entity with a new ID and CreatedAt if necessary. Updates UpdatedAt.
func (e Metadata) Clone(now time.Time) Metadata {
	return Metadata{
		ID:        e.ID.SelfOrNew(),
		CreatedAt: e.createdAtOrNow(now),
		UpdatedAt: now,
		Version:   e.Version + 1,
		TenantID:  e.TenantID,
	}
}

func (e Metadata) CloneIn(ctx context.Context, now time.Time) Metadata {
	clone := e.Clone(now)
	if e.TenantID == "" {
		clone.TenantID = xcontext.TenantOrDefault(ctx)
	}
	return clone
}

func (e Metadata) createdAtOrNow(now time.Time) time.Time {
	var created = e.CreatedAt
	if !created.IsZero() {
		return created
	}
	return now
}
