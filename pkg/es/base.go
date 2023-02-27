package es

import (
	"github.com/AltScore/gothic/pkg/ids"
	"time"
)

type Metadata struct {
	entityType string
	entityID   ids.ID
	type_      string
	version    int
	datetime   time.Time
}

func NewMetadata(
	entityType string,
	entityID ids.ID,
	type_ string,
	version int) Metadata {
	return Metadata{
		entityType: entityType,
		entityID:   entityID,
		type_:      type_,
		version:    version,
		datetime:   time.Now(),
	}
}

func (m *Metadata) EntityType() string {
	return m.entityType
}

func (m *Metadata) EntityID() ids.ID {
	return m.entityID
}

func (m *Metadata) Type() string {
	return m.type_
}

func (m *Metadata) Version() int {
	return m.version
}

func (m *Metadata) Datetime() time.Time {
	return m.datetime
}
