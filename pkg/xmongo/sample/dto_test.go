package sample

import (
	"github.com/AltScore/gothic/pkg/xmongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"testing"
)

func TestCanEncodeYDecode(t *testing.T) {
	p := NewBuilder().
		WithId("1").
		WithName("John").
		WithAge(30).
		Build()

	r := makeRegistry()

	// WHEN marshall
	personBson, err := bson.MarshalWithRegistry(r, p)

	require.NoError(t, err)

	// THEN can unmarshall
	var actualPerson Person
	err = bson.UnmarshalWithRegistry(r, personBson, &actualPerson)

	require.NoError(t, err)

	assert.Equal(t, p, actualPerson)
}

func makeRegistry() *bsoncodec.Registry {
	rb := bson.NewRegistryBuilder()

	de := xmongo.NewDecoderEncoder[Person, personDto, person](
		fromEntity,
		func(p personDto) Person { return p.ToEntity() },
	)

	de.Register(rb)

	return rb.Build()
}
