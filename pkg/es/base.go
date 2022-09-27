package es

import (
	"time"
)

type Metadata[ID EntityID] struct {
	entityType string
	entityID   ID
	type_      string
	version    int
	datetime   time.Time
}

func NewMetadata[ID EntityID](
	entityType string,
	entityID ID,
	type_ string,
	version int) Metadata[ID] {
	return Metadata[ID]{
		entityType: entityType,
		entityID:   entityID,
		type_:      type_,
		version:    version,
		datetime:   time.Now(),
	}
}

func (b Metadata[ID]) EntityType() string {
	return b.entityType
}

func (b Metadata[ID]) EntityID() ID {
	return b.entityID
}

func (b Metadata[ID]) Type() string {
	return b.type_
}

func (b Metadata[ID]) Version() int {
	return b.version
}

func (b Metadata[ID]) Datetime() time.Time {
	return b.datetime
}
