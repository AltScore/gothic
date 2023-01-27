package event

import (
	"reflect"
	"time"

	"github.com/AltScore/gothic/pkg/es/codec"
	"go.mongodb.org/mongo-driver/bson"
)

type entry struct {
	ID               string      `bson:"_id"`
	Name             string      `bson:"name"`
	Time             time.Time   `bson:"time"`
	TimeNano         int64       `bson:"timeNano"`
	AggregateName    string      `bson:"aggName"`
	AggregateID      AggID       `bson:"aggID"`
	AggregateVersion int         `bson:"aggV"`
	Data             interface{} `bson:"data"`
}

// Intermediary struct to hold the data from the bson without knowing the type
type rawEntry struct {
	ID               string    `bson:"_id"`
	Name             string    `bson:"name"`
	Time             time.Time `bson:"time"`
	TimeNano         int64     `bson:"timeNano"`
	AggregateName    string    `bson:"aggName"`
	AggregateID      AggID     `bson:"aggID"`
	AggregateVersion int       `bson:"aggV"`
	Data             bson.Raw  `bson:"data"`
}

func From(event Event) *entry {
	id, name, version := event.Aggregate()

	return &entry{
		ID:               event.ID().String(),
		Name:             event.Name(),
		Time:             event.Time(),
		TimeNano:         event.Time().UnixNano(),
		AggregateName:    name,
		AggregateID:      id,
		AggregateVersion: version,
		Data:             event.Data(),
	}
}

func (e *entry) Metadata() *Metadata {
	return &Metadata{
		ID:               ID(e.ID),
		Name:             e.Name,
		Time:             time.Unix(e.Time.Unix(), e.TimeNano),
		AggregateName:    e.AggregateName,
		AggregateID:      e.AggregateID,
		AggregateVersion: e.AggregateVersion,
		Data:             e.Data,
	}
}

func (e *entry) ToEvent() (Event, error) {
	return Event{m: e.Metadata()}, nil
}

func (e *entry) UnmarshalBSON(data []byte) error {
	eX := rawEntry{}

	if err := bson.Unmarshal(data, &eX); err != nil {
		return err
	}

	e.ID = eX.ID
	e.Name = eX.Name
	e.Time = eX.Time
	e.TimeNano = eX.TimeNano
	e.AggregateName = eX.AggregateName
	e.AggregateID = eX.AggregateID
	e.AggregateVersion = eX.AggregateVersion

	if eX.Data == nil {
		return nil
	}

	factory, found := codec.ForName(e.Name)

	if !found {
		panic("no codec for " + e.Name)
	}

	var err error
	v := reflect.ValueOf(&e.Data)
	if err = factory(eX.Data, v); err != nil {
		return err
	}

	return nil
}

func (e Event) MarshalBSON() ([]byte, error) {
	return bson.Marshal(From(e))
}

func (e *Event) UnmarshalBSON(data []byte) error {
	var en entry

	if err := bson.Unmarshal(data, &en); err != nil {
		return err
	}

	e.m = en.Metadata()
	return nil
}
