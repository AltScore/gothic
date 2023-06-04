package xbson

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v3"
)

type sampleType string

func (s sampleType) String() string {
	return string(s)
}

var _ fmt.Stringer = sampleType("")

type sampleTypedGeneric interface {
	Typed
}

type sampleTypedOne struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func (s sampleTypedOne) T() string { return "sample-type-one" }

type sampleTypedTwo struct {
	Email   string `bson:"email"`
	Enabled bool   `bson:"enabled"`
}

func (s sampleTypedTwo) T() string { return "sample-type-two" }

var _ sampleTypedGeneric = &sampleTypedOne{}

func Test_TypedGeneric_can_encode_and_decode(t *testing.T) {

	// Given a bson registry
	builder := bson.NewRegistryBuilder()
	sampleCodex := NewTypedGenericCodex[sampleTypedGeneric]()
	sampleCodex.RegisterType(func() Typed { return &sampleTypedOne{} })

	sampleCodex.Register(builder)

	registry := builder.Build()

	// Given a typed generic
	var expected sampleTypedGeneric = &sampleTypedOne{Name: "Alice", Age: 24}

	// When it is encoded
	bytes, err := bson.MarshalWithRegistry(registry, &expected)

	dumpBson(t, "-- Sample value", bytes)

	// Then it should not fail
	require.NoError(t, err)

	hexBytes := hex.EncodeToString(bytes)
	fmt.Println(hexBytes)

	// When it is decoded
	var actual sampleTypedGeneric
	err = bson.UnmarshalWithRegistry(registry, bytes, &actual)

	// Then it should not fail
	require.NoError(t, err)

	// And it should be the same as the original
	require.Equal(t, expected, actual)
}

type sampleStruct struct {
	Generic sampleTypedGeneric `bson:"generic"`
}

func Test_TypedGeneric_can_be_embedded(t *testing.T) {

	// Given a bson registry
	builder := bson.NewRegistryBuilder()
	sampleCodex := NewTypedGenericCodex[sampleTypedGeneric]()
	sampleCodex.RegisterType(func() Typed { return &sampleTypedOne{} })

	sampleCodex.Register(builder)

	registry := builder.Build()

	// Given a typed generic
	var expected = sampleStruct{
		Generic: &sampleTypedOne{Name: "Alice", Age: 24},
	}

	// When it is encoded
	bytes, err := bson.MarshalWithRegistry(registry, &expected)

	dumpBson(t, "-- Sample value", bytes)

	// Then it should not fail
	require.NoError(t, err)

	hexBytes := hex.EncodeToString(bytes)
	fmt.Println(hexBytes)

	// When it is decoded
	var actual sampleStruct
	err = bson.UnmarshalWithRegistry(registry, bytes, &actual)

	// Then it should not fail
	require.NoError(t, err)

	// And it should be the same as the original
	require.Equal(t, expected, actual)
}

type strucWithSlice struct {
	Generic []sampleTypedGeneric `bson:"generic"`
}

func Test_TypedGeneric_can_be_embedded_in_a_slice(t *testing.T) {

	// Given a bson registry
	builder := bson.NewRegistryBuilder()
	sampleCodex := NewTypedGenericCodex[sampleTypedGeneric]()
	sampleCodex.RegisterType(func() Typed { return &sampleTypedOne{} })
	sampleCodex.RegisterType(func() Typed { return &sampleTypedTwo{} })

	sampleCodex.Register(builder)

	registry := builder.Build()

	// Given a typed generic
	var expected = strucWithSlice{
		Generic: []sampleTypedGeneric{
			&sampleTypedOne{Name: "Alice", Age: 24},
			&sampleTypedTwo{Email: "sample@s.com", Enabled: true},
		},
	}

	// When it is encoded
	bytes, err := bson.MarshalWithRegistry(registry, &expected)

	dumpBson(t, "-- Sample value", bytes)

	// Then it should not fail
	require.NoError(t, err)

	hexBytes := hex.EncodeToString(bytes)
	fmt.Println(hexBytes)

	// When it is decoded
	var actual strucWithSlice
	err = bson.UnmarshalWithRegistry(registry, bytes, &actual)

	// Then it should not fail
	require.NoError(t, err)

	// And it should be the same as the original
	require.Equal(t, expected, actual)
}

func dumpBson(t *testing.T, title string, bytes []byte) {
	var bsonMap bson.M

	err := bson.Unmarshal(bytes, &bsonMap)
	require.NoError(t, err)

	// Convert object to YAML
	yamlBytes, err := yaml.Marshal(bsonMap)

	require.NoError(t, err)

	yamlString := string(yamlBytes)
	fmt.Println(title)
	fmt.Println(yamlString)

	fmt.Println(bsonMap)
}
