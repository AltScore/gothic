package event

import (
	"fmt"
	"testing"

	"github.com/AltScore/gothic/pkg/es/codec"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	EmptyEventName         = "empty.type"
	EventWithDataName      = "withData.type"
	EventWithOtherDataName = "other.type"
)

type emptyEventData struct {
}

type eventWithData struct {
	Str string `bson:"str"`
	Num int    `bson:"num"`
}

type eventWithOtherData struct {
	Num     float64 `bson:"num"`
	YesOrNo bool    `bson:"yesOrNo"`
}

func Test_EmptyEvent_can_be_marshalled(t *testing.T) {
	event := New(EmptyEventName, emptyEventData{}, WithID("sample-id"))

	// WHEN we marshal it
	bsonBytes, err := marshal(event)

	// THEN we expect no error
	require.NoError(t, err)

	// THEN we expect the bsonBytes to be correct
	var actual map[string]interface{}

	err = bson.Unmarshal(bsonBytes, &actual)

	fmt.Println(actual)

	require.NoError(t, err)
	require.Equal(t, "sample-id", actual["_id"])
	require.Equal(t, EmptyEventName, actual["name"])
	require.Equal(t, int32(0), actual["aggV"])
	require.Equal(t, "", actual["aggName"])
	require.Equal(t, "", actual["aggID"])
}

func Test_EventWithData_can_be_marshalled(t *testing.T) {
	data := eventWithData{
		Str: "a nice value",
		Num: 42,
	}
	event := New(EventWithDataName, data, WithID("another-id"))

	// WHEN we marshal it
	bsonBytes, err := marshal(event)

	// THEN we expect no error
	require.NoError(t, err)

	// THEN we expect the bsonBytes to be correct
	var actual map[string]interface{}

	err = bson.Unmarshal(bsonBytes, &actual)

	fmt.Println(actual)

	require.NoError(t, err)
	require.Equal(t, "another-id", actual["_id"])
	require.Equal(t, EventWithDataName, actual["name"])
	require.Equal(t, int32(0), actual["aggV"])
	require.Equal(t, "", actual["aggName"])
	require.Equal(t, "", actual["aggID"])

	actualDataMap := actual["data"]

	require.NotNil(t, actualDataMap)

	dataData := actualDataMap.(map[string]interface{})

	require.Equal(t, "a nice value", dataData["str"])
	require.Equal(t, int32(42), dataData["num"])
}

func Test_EmptyEvent_can_be_unmarshalled(t *testing.T) {
	codec.Register[emptyEventData](EmptyEventName)

	// GIVEN a mashalled empty event
	event := New(EmptyEventName, emptyEventData{}, WithID("sample-id"))
	bsonBytes, err := marshal(event)
	require.NoError(t, err)

	// WHEN we unmarshal it
	actual, err := unmarshal(bsonBytes)

	// THEN we expect no error
	require.NoError(t, err)

	// THEN we expect the event to be correct
	id, name, version := actual.Aggregate()

	require.Equal(t, "sample-id", actual.ID())
	require.Equal(t, EmptyEventName, actual.Name())
	require.Equal(t, 0, version)
	require.Equal(t, "", name)
	require.Equal(t, "", id)
}

func Test_EventWithData_can_be_unmarshalled(t *testing.T) {
	codec.Register[eventWithData](EventWithDataName)

	// GIVEN a mashalled empty event
	data := eventWithData{
		Str: "a nice value",
		Num: 42,
	}
	event := New(EventWithDataName, data, WithID("another-id"))

	bsonBytes, err := marshal(event)
	require.NoError(t, err)

	// WHEN we unmarshal it
	actual, err := unmarshal(bsonBytes)

	// THEN we expect no error
	require.NoError(t, err)

	// THEN we expect the event to be correct
	id, name, version := actual.Aggregate()

	d := actual.Data()
	actualData, ok := d.(*eventWithData)

	require.True(t, ok)

	require.Equal(t, "another-id", actual.ID())
	require.Equal(t, EventWithDataName, actual.Name())
	require.Equal(t, 0, version)
	require.Equal(t, "", name)
	require.Equal(t, "", id)

	require.Equal(t, &data, actualData)
}

func Test_Can_marshall_list_of_events(t *testing.T) {
	codec.Register[emptyEventData](EmptyEventName)
	codec.Register[eventWithData](EventWithDataName)
	codec.Register[eventWithOtherData](EventWithOtherDataName)

	// GIVEN a list of events
	events := []Event{
		New(EmptyEventName, emptyEventData{}, WithID("sample-id")),
		New(EventWithDataName, eventWithData{
			Str: "a nice value",
			Num: 42,
		}, WithID("another-id")),
		New(EventWithOtherDataName, eventWithOtherData{
			Num:     3.14,
			YesOrNo: true,
		}, WithID("yet-another-id")),
	}

	// WHEN we marshal it
	bsonBytes, err := marshalAll(events)

	// THEN we expect no error
	require.NoError(t, err)

	// THEN we expect the bsonBytes to be correct
	actualEvents := unmarshalAll(bsonBytes)

	require.Len(t, actualEvents, 3)
	require.Equal(t, "sample-id", actualEvents[0].ID())
	require.Equal(t, EmptyEventName, actualEvents[0].Name())

	require.Equal(t, "another-id", actualEvents[1].ID())
	require.Equal(t, EventWithDataName, actualEvents[1].Name())
	require.Equal(t, "a nice value", actualEvents[1].Data().(*eventWithData).Str)

	require.Equal(t, "yet-another-id", actualEvents[2].ID())
	require.Equal(t, EventWithOtherDataName, actualEvents[2].Name())
	require.Equal(t, 3.14, actualEvents[2].Data().(*eventWithOtherData).Num)
}

func TestCan_marshall_and_unmarshall_event_directly(t *testing.T) {
	codec.Register[eventWithOtherData](EventWithOtherDataName)

	// GIVEN an event
	event := New(EventWithOtherDataName, eventWithOtherData{
		Num:     3.14,
		YesOrNo: true,
	}, WithID("yet-another-id"))

	// WHEN we marshal it
	bsonBytes, err := bson.Marshal(event)

	// THEN we expect no error
	require.NoError(t, err)

	// WHEN we unmarshal it
	var actual Event
	err = bson.Unmarshal(bsonBytes, &actual)

	// THEN we expect no error
	require.NoError(t, err)

	// THEN we expect the event to be correct
	require.Equal(t, "yet-another-id", actual.ID())
	require.Equal(t, EventWithOtherDataName, actual.Name())
}

type docWithEvents struct {
	Entries []entry `bson:"events"`
}

func unmarshalAll(bytes []byte) []Event {

	doc := docWithEvents{}

	err := bson.Unmarshal(bytes, &doc)

	if err != nil {
		panic(err)
	}
	var events []Event
	for _, entry := range doc.Entries {
		event, err := entry.ToEvent()
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	return events
}

func marshalAll(events []Event) ([]byte, error) {
	var bsonEvents []entry

	for _, event := range events {
		bsonEvents = append(bsonEvents, *From(event))
	}

	doc := docWithEvents{
		Entries: bsonEvents,
	}
	return bson.Marshal(doc)
}

func marshal(e Event) ([]byte, error) {
	entry := From(e)

	return bson.Marshal(entry)
}

func unmarshal(b []byte) (Event, error) {
	var entry entry
	err := bson.Unmarshal(b, &entry)
	if err != nil {
		return EmptyEvent, err
	}

	return entry.ToEvent()
}
