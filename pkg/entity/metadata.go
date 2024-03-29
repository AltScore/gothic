package entity

import (
	"context"
	"github.com/AltScore/gothic/v2/pkg/ids"
	"github.com/AltScore/gothic/v2/pkg/xcontext"
	"github.com/AltScore/gothic/v2/pkg/xeh"
	eh "github.com/looplab/eventhorizon"
	"time"
)

// Metadata is the metadata for any entity.
type Metadata struct {
	ID        ids.Id    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Version   int       `json:"version" bson:"version"`
	Tenant    string    `json:"tenant,omitempty" bson:"tenant,omitempty"`
}

type Option func(m *Metadata)

func New(options ...Option) Metadata {
	m := Metadata{
		ID: ids.New(),
	}

	for _, option := range options {
		option(&m)
	}
	return m
}

// FromEvent creates an entity metadata from the EventHorizon event
func FromEvent(event eh.Event) Metadata {
	return Metadata{
		ID:        ids.OrNew(event.AggregateID()),
		CreatedAt: event.Timestamp(),
		UpdatedAt: event.Timestamp(),
		Version:   event.Version(),
		Tenant:    xeh.GetEventTenant(event),
	}
}

func WithId(id ids.Id) Option {
	return func(m *Metadata) {
		m.ID = id
	}
}

func At(t time.Time) Option {
	return func(m *Metadata) {
		m.CreatedAt = t
		m.UpdatedAt = t
	}
}

func WithTenant(tenant string) Option {
	return func(m *Metadata) {
		m.Tenant = tenant
	}
}

func WithCtx(ctx context.Context) Option {
	return func(m *Metadata) {
		m.Tenant = xcontext.GetTenantOrDefault(ctx)
	}
}

// GetID returns the ID of the entity. Implements Identifiable interfaces.
func (e Metadata) GetID() ids.Id {
	return e.ID
}

func (e Metadata) GetTenant() string {
	return e.Tenant
}

// Clone returns a clone of the entity with a new ID and CreatedAt if necessary. Updates UpdatedAt.
func (e Metadata) Clone(now time.Time) Metadata {
	return Metadata{
		ID:        ids.OrNew(e.ID),
		CreatedAt: e.createdAtOrNow(now),
		UpdatedAt: now,
		Version:   e.Version + 1,
		Tenant:    e.Tenant,
	}
}

func (e Metadata) CloneIn(ctx context.Context, now time.Time) Metadata {
	clone := e.Clone(now)
	if e.Tenant == "" {
		clone.Tenant = xcontext.GetTenantOrDefault(ctx)
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
