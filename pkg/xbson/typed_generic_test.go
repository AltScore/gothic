package xbson

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

type sampleType string

func (s sampleType) String() string {
	return string(s)
}

var _ fmt.Stringer = sampleType("")

type sampleTypedGeneric interface {
	T() string
}

type sampleTypedOne struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func (s sampleTypedOne) T() string { return "sample-type-one" }

type sampleTypedOneDto sampleTypedOne

type sampleTypedTwo struct {
	Email   string `bson:"email"`
	Enabled bool   `bson:"enabled"`
}

func (s sampleTypedTwo) T() string { return "sample-type-two" }

var _ sampleTypedGeneric = &sampleTypedOne{}

type sampleTypedTwoDto sampleTypedTwo

func Test_TypedGeneric_can_encode_and_decode(t *testing.T) {

	// Given a bson registry
	registry := bson.NewRegistry()

	sampleCodex := NewTypedGenericCodex[sampleTypedGeneric](
		func(t sampleTypedGeneric) string { return t.T() },
	)

	sampleCodex.RegisterType(
		func() sampleTypedGeneric { return &sampleTypedOne{} },
		func(t sampleTypedGeneric) interface{} { return (*sampleTypedOneDto)(t.(*sampleTypedOne)) },
		func(dto interface{}) sampleTypedGeneric { return (*sampleTypedOne)(dto.(*sampleTypedOneDto)) },
	)

	sampleCodex.Register(registry)

	// Given a typed generic
	var expected sampleTypedGeneric = &sampleTypedOne{Name: "Alice", Age: 24}

	// When it is encoded
	bs, err := MarshalWithRegistry(registry, &expected)

	Dump(t, "-- Sample value", bs)

	// Then it should not fail
	require.NoError(t, err)

	hexBytes := hex.EncodeToString(bs)
	fmt.Println(hexBytes)

	// When it is decoded
	var actual sampleTypedGeneric
	err = UnmarshalWithRegistry(registry, bs, &actual)

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
	registry := bson.NewRegistry()

	sampleCodex := NewTypedGenericCodex[sampleTypedGeneric](
		func(t sampleTypedGeneric) string { return t.T() },
	)
	sampleCodex.RegisterType(
		func() sampleTypedGeneric { return &sampleTypedOne{} },
		func(t sampleTypedGeneric) interface{} { return (*sampleTypedOneDto)(t.(*sampleTypedOne)) },
		func(dto interface{}) sampleTypedGeneric { return (*sampleTypedOne)(dto.(*sampleTypedOneDto)) },
	)

	sampleCodex.Register(registry)

	// Given a typed generic
	var expected = sampleStruct{
		Generic: &sampleTypedOne{Name: "Alice", Age: 24},
	}

	// When it is encoded
	bs, err := MarshalWithRegistry(registry, &expected)

	Dump(t, "-- Sample value", bs)

	// Then it should not fail
	require.NoError(t, err)

	hexBytes := hex.EncodeToString(bs)
	fmt.Println(hexBytes)

	// When it is decoded
	var actual sampleStruct
	err = UnmarshalWithRegistry(registry, bs, &actual)

	// Then it should not fail
	require.NoError(t, err)

	// And it should be the same as the original
	require.Equal(t, expected, actual)
}

func UnmarshalWithRegistry(registry *bsoncodec.Registry, bs []byte, value interface{}) error {
	hexBytes := hex.EncodeToString(bs)
	fmt.Println(hexBytes)

	dec, err := bson.NewDecoder(bsonrw.NewBSONDocumentReader(bs))
	if err != nil {
		return err
	}

	if err := dec.SetRegistry(registry); err != nil {
		return err
	}

	return dec.Decode(value)
}

func MarshalWithRegistry(registry *bsoncodec.Registry, value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	vw, err := bsonrw.NewBSONValueWriter(buf)
	if err != nil {
		panic(err)
	}
	enc, err := bson.NewEncoder(vw)
	if err != nil {
		return nil, err
	}

	if err := enc.SetRegistry(registry); err != nil {
		return nil, err
	}

	if err := enc.Encode(value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type strucWithSlice struct {
	Generic []sampleTypedGeneric `bson:"generic"`
}

func Test_TypedGeneric_can_be_embedded_in_a_slice(t *testing.T) {

	// Given a bson registry
	registry := bson.NewRegistry()
	sampleCodex := NewTypedGenericCodex[sampleTypedGeneric](
		func(t sampleTypedGeneric) string { return t.T() },
	)
	sampleCodex.RegisterType(func() sampleTypedGeneric { return &sampleTypedOne{} },
		func(t sampleTypedGeneric) interface{} { return (*sampleTypedOneDto)(t.(*sampleTypedOne)) },
		func(dto interface{}) sampleTypedGeneric { return (*sampleTypedOne)(dto.(*sampleTypedOneDto)) },
	)
	sampleCodex.RegisterType(func() sampleTypedGeneric { return &sampleTypedTwo{} },
		func(t sampleTypedGeneric) interface{} { return (*sampleTypedTwoDto)(t.(*sampleTypedTwo)) },
		func(dto interface{}) sampleTypedGeneric { return (*sampleTypedTwo)(dto.(*sampleTypedTwoDto)) },
	)

	sampleCodex.Register(registry)

	// Given a typed generic
	var expected = strucWithSlice{
		Generic: []sampleTypedGeneric{
			&sampleTypedOne{Name: "Alice", Age: 24},
			&sampleTypedTwo{Email: "sample@s.com", Enabled: true},
		},
	}

	// When it is encoded
	bs, err := MarshalWithRegistry(registry, &expected)

	Dump(t, "-- Sample value", bs)

	// Then it should not fail
	require.NoError(t, err)

	hexBytes := hex.EncodeToString(bs)
	fmt.Println(hexBytes)

	// When it is decoded
	var actual strucWithSlice
	err = UnmarshalWithRegistry(registry, bs, &actual)

	// Then it should not fail
	require.NoError(t, err)

	// And it should be the same as the original
	require.Equal(t, expected, actual)
}

func Test_Managing_reflect_values(t *testing.T) {

	valueOfTypeOne := sampleTypedOne{
		Name: "Morgana",
		Age:  42,
	}

	pointerToValueOfTypeOne := &valueOfTypeOne

	var sampleTypedGenericOfValue sampleTypedGeneric
	sampleTypedGenericOfValue = valueOfTypeOne

	var sampleTypedGenericOfPointer sampleTypedGeneric
	sampleTypedGenericOfPointer = pointerToValueOfTypeOne

	// Now, the type
	reflectValueOfTypeOne := reflect.ValueOf(valueOfTypeOne)
	reflectedPointerToValueOfTypeOne := reflect.ValueOf(pointerToValueOfTypeOne)
	reflectedSampleTypedGenericOfValue := reflect.ValueOf(sampleTypedGenericOfValue)
	reflectedSampleTypedGenericOfPointer := reflect.ValueOf(sampleTypedGenericOfPointer)

	// And the type of the types
	assert.Equal(t, reflect.Struct, reflectValueOfTypeOne.Kind())
	assert.Equal(t, reflect.Ptr, reflectedPointerToValueOfTypeOne.Kind())
	assert.Equal(t, reflect.Struct, reflectedSampleTypedGenericOfValue.Kind())
	assert.Equal(t, reflect.Ptr, reflectedSampleTypedGenericOfPointer.Kind())

	// And the type of the types of the types
	assert.Equal(t, reflect.Struct, reflectValueOfTypeOne.Type().Kind())
	assert.Equal(t, reflect.Struct, reflectedSampleTypedGenericOfValue.Type().Kind())

	assert.Equal(t, reflect.Ptr, reflectedPointerToValueOfTypeOne.Type().Kind())
	assert.Equal(t, reflect.Ptr, reflectedSampleTypedGenericOfPointer.Type().Kind())

	// And the names of the types

	assert.Equal(t, "sampleTypedOne", reflectValueOfTypeOne.Type().Name())
	assert.Equal(t, "sampleTypedOne", reflectedPointerToValueOfTypeOne.Type().Elem().Name())
	assert.Equal(t, "sampleTypedOne", reflectedSampleTypedGenericOfValue.Type().Name())
	assert.Equal(t, "sampleTypedOne", reflectedSampleTypedGenericOfPointer.Type().Elem().Name())

	//
	elementOfInterface := reflectedSampleTypedGenericOfPointer.Elem().Interface()
	assert.Equal(t, "sampleTypedOne", reflect.TypeOf(elementOfInterface).Name())

	//
	original := reflectedSampleTypedGenericOfPointer.Interface().(interface{})

	reflectOriginal := reflect.ValueOf(original)
	assert.Equal(t, "sampleTypedOne", reflectOriginal.Type().Elem().Name())

	// Types are the same, there is no runtime info about interfaces in the values
	assert.Equal(t, reflectOriginal.Type(), reflectedSampleTypedGenericOfPointer.Type())

}
