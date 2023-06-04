package xbson

import (
	"fmt"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

// Typed is implemented by a family of types that can be encoded/decoded to/from a bson document.
// It allows a generic type to be encoded/decoded to/from a bson document.
type Typed interface {
	// T returns the type of the document. Different types should report different values.
	T() string
}

// TypedGenericCodex is a generic encoder/decoder for a family of types that implement the Typed interface
// It allows a generic type to be encoded/decoded to/from a bson document.
// The method T() is used to determine the type of the document.
type TypedGenericCodex struct {
	factoryByType map[string]func() Typed
	valueType     reflect.Type
	lock          sync.RWMutex
}

var _ EncoderDecoder = (*TypedGenericCodex)(nil)

func NewTypedGenericCodex[Interface Typed]() *TypedGenericCodex {

	return &TypedGenericCodex{
		factoryByType: make(map[string]func() Typed),
		valueType:     reflect.TypeOf((*Interface)(nil)).Elem(), // The type of the interface
	}
}

var _ bsoncodec.ValueCodec = (*TypedGenericCodex)(nil)
var _ bsoncodec.ValueEncoder = (*TypedGenericCodex)(nil)

// Register implements the bsoncodec.RegistryBuilder interface
// It allows the decoderEncoder to be registered with a bsoncodec.RegistryBuilder
func (t *TypedGenericCodex) Register(builder *bsoncodec.RegistryBuilder) {

	builder.RegisterTypeDecoder(t.valueType, t)
	builder.RegisterTypeEncoder(t.valueType, t)
}

type wrapper struct {
	T string
	V bson.Raw
}

// RegisterType registers a factory function for a given type name
func (t *TypedGenericCodex) RegisterType(factory func() Typed) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.factoryByType[factory().T()] = factory
}

// LookupType returns the factory function for a given type name
func (t *TypedGenericCodex) LookupType(typeName string) (func() Typed, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	factory, found := t.factoryByType[typeName]
	return factory, found
}

// EncodeValue implements the bsoncodec.ValueEncoder interface
func (t *TypedGenericCodex) EncodeValue(ctx bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
	// Encode the original underlying value (it is the struct, not the interface)
	bytes, err := bson.Marshal(value.Interface())

	if err != nil {
		return err
	}

	// Wrap the original value with its type
	v := wrapper{T: value.Interface().(Typed).T(), V: bytes}

	// Encode the wrapped value
	encoder, err := ctx.Registry.LookupEncoder(reflect.TypeOf(v))
	if err != nil {
		return err
	}

	return encoder.EncodeValue(ctx, writer, reflect.ValueOf(v))
}

// DecodeValue implements the bsoncodec.ValueDecoder interface
func (t *TypedGenericCodex) DecodeValue(ctx bsoncodec.DecodeContext, reader bsonrw.ValueReader, value reflect.Value) error {
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
