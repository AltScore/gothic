package es

import "time"

type Event[ID EntityID[ID], Snapshot any] interface {
	EntityType() string
	EntityID() ID
	Type() string
	Version() int
	Datetime() time.Time
	Apply(*Snapshot) error
}
