package xbson

import (
	"fmt"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

// GetType provides the type for a family of types that can be encoded/decoded to/from a bson document.
// It allows a generic type to be encoded/decoded to/from a bson document.
type GetType[Typed any] func(Typed) string

// TypedGenericCodex is a generic encoder/decoder for a family of types that implement the Typed interface
// It allows a generic type to be encoded/decoded to/from a bson document.
// The method T() is used to det/ermine the type of the document.
type TypedGenericCodex[Typed any] struct {
	factoryByType map[string]func() Typed
	toDtoByType   map[string]func(Typed) interface{}
	getType       GetType[Typed]
	valueType     reflect.Type
	lock          sync.RWMutex
}

var _ EncoderDecoder = (*TypedGenericCodex[string])(nil)

func NewTypedGenericCodex[Typed any](getType GetType[Typed]) *TypedGenericCodex[Typed] {
	return &TypedGenericCodex[Typed]{
		factoryByType: make(map[string]func() Typed),
		toDtoByType:   make(map[string]func(Typed) interface{}),
		getType:       getType,
		valueType:     reflect.TypeOf((*Typed)(nil)).Elem(), // The type of the interface
	}
}

var _ bsoncodec.ValueCodec = (*TypedGenericCodex[string])(nil)
var _ bsoncodec.ValueEncoder = (*TypedGenericCodex[string])(nil)

// Register implements the bsoncodec.RegistryBuilder interface
// It allows the decoderEncoder to be registered with a bsoncodec.RegistryBuilder
func (t *TypedGenericCodex[Typed]) Register(builder *bsoncodec.RegistryBuilder) {

	builder.RegisterHookEncoder(t.valueType, t)
	builder.RegisterHookDecoder(t.valueType, t)
}

type wrapper struct {
	T string
	V bson.Raw
}

// RegisterType registers a factory function for a given type name
func (t *TypedGenericCodex[Typed]) RegisterType(factory func() Typed, toDto func(Typed) interface{}) {
	t.lock.Lock()
	defer t.lock.Unlock()

	type_ := t.getType(factory())
	t.factoryByType[type_] = factory
	t.toDtoByType[type_] = toDto
}

// LookupType returns the factory function for a given type name
func (t *TypedGenericCodex[Typed]) LookupType(typeName string) (func() Typed, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	factory, found := t.factoryByType[typeName]
	return factory, found
}

// EncodeValue implements the bsoncodec.ValueEncoder interface
func (t *TypedGenericCodex[Typed]) EncodeValue(ctx bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
	// Encode the original underlying value (it is the struct, not the interface)
	typed, ok := value.Interface().(Typed)
	if !ok {
		return fmt.Errorf("value does not implement Typed interface")
	}

	typeName := t.getType(typed)

	toDto := t.toDtoByType[typeName]
	dto := toDto(typed)
	bytes, err := bson.MarshalWithRegistry(ctx.Registry, dto) // use same registry

	if err != nil {
		return err
	}

	// Wrap the original value with its type
	v := wrapper{T: typeName, V: bytes}

	// Encode the wrapped value
	encoder, err := ctx.Registry.LookupEncoder(reflect.TypeOf(v))
	if err != nil {
		return err
	}

	return encoder.EncodeValue(ctx, writer, reflect.ValueOf(v))
}

// DecodeValue implements the bsoncodec.ValueDecoder interface
func (t *TypedGenericCodex[Typed]) DecodeValue(ctx bsoncodec.DecodeContext, reader bsonrw.ValueReader, value reflect.Value) error {
	// Decode the wrapped value
	var v wrapper
	decoder, err := ctx.Registry.LookupDecoder(reflect.TypeOf(&v).Elem())
	if err != nil {
		return err
	}

	err = decoder.DecodeValue(ctx, reader, reflect.ValueOf(&v).Elem())
	if err != nil {
		return err
	}

	// Create an instance of the original underlying value type
	factory, found := t.LookupType(v.T)

	if !found {
		return fmt.Errorf("unknown type: %s", v.T)
	}

	// Decode the original underlying value
	result := factory()
	err = bson.Unmarshal(v.V, result)
	if err != nil {
		return err
	}

	// Set the original underlying value to the target interface
	value.Set(reflect.ValueOf(result))

	return nil
}
