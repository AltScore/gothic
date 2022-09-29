package event

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

type ID = string
type AggID = string

type Aggregate interface {
	ID() string
	Type() string
	Version() int
}

type IEvent interface {
	// ID returns the id of the event.
	ID() ID
	// Name returns the name of the event.
	Name() string
	// Time returns the time of the event.
	Time() time.Time
	// Data returns the event data.
	Data() interface{}
	// Aggregate returns the id, name and version of the aggregate that the
	// event belongs to. aggregate should return zero values if the event is not
	// an aggregate event.
	Aggregate() (id string, name string, version int)
}

type Metadata struct {
	ID               ID
	Name             string
	Time             time.Time
	AggregateName    string
	AggregateID      AggID
	AggregateVersion int
	Data             interface{}
}

type Event struct {
	m *Metadata
}

func (e Event) Data() interface{} {
	return e.m.Data
}

type Option func(*Metadata)

func New(name string, data any, opts ...Option) IEvent {
	m := Metadata{
		ID:   uuid.New().String(),
		Name: name,
		Time: time.Now(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	return Event{m: &Metadata{
		ID:               m.ID,
		Name:             m.Name,
		Time:             m.Time,
		AggregateName:    m.AggregateName,
		AggregateID:      m.AggregateID,
		AggregateVersion: m.AggregateVersion,
		Data:             data,
	}}
}

func For[Data any](a Aggregate, name string, data Data, opts ...Option) IEvent {
	return New(name, data, append(opts, WithAggregate(a))...)
}

func (e Event) ID() ID {
	return e.m.ID
}

func (e Event) Name() string {
	return e.m.Name
}

func (e Event) Time() time.Time {
	return e.m.Time
}

func (e Event) Aggregate() (id, name string, version int) {
	return e.m.AggregateID, e.m.AggregateName, e.m.AggregateVersion
}

func WithID(id ID) Option {
	return func(m *Metadata) {
		m.ID = id
	}
}

func WithTime(t time.Time) Option {
	return func(m *Metadata) {
		m.Time = t
	}
}

func WithAggregate(a Aggregate) Option {
	return func(m *Metadata) {
		m.AggregateID = a.ID()
		m.AggregateName = a.Type()
		m.AggregateVersion = a.Version() + 1
	}
}

func DataOf[T any](e IEvent) T {
	if data, ok := e.Data().(T); ok {
		return data
	}
	var t T

	panic("Event: data is not of type " + reflect.TypeOf(t).String())
}
