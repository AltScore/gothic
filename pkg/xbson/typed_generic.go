package xbson

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// GetType provides the type for a family of types that can be encoded/decoded to/from a bson document.
// It allows a generic type to be encoded/decoded to/from a bson document.
type GetType[Typed any] func(Typed) string

type subtype[Typed any] struct {
	factory func() Typed
	toDto   func(Typed) interface{}
	fromDto func(interface{}) Typed
}

// TypedGenericCodex is a generic encoder/decoder for a family of types that implement the Typed interface
// It allows a generic type to be encoded/decoded to/from a bson document.
// The getType function is used to determine the type of the underlying value.
type TypedGenericCodex[Typed any] struct {
	subtypes  map[string]subtype[Typed]
	getType   GetType[Typed]
	valueType reflect.Type
	lock      sync.RWMutex
}

var _ EncoderDecoder = (*TypedGenericCodex[string])(nil)

func NewTypedGenericCodex[Typed any](getType GetType[Typed]) *TypedGenericCodex[Typed] {
	return &TypedGenericCodex[Typed]{
		subtypes:  make(map[string]subtype[Typed]),
		getType:   getType,
		valueType: reflect.TypeOf((*Typed)(nil)).Elem(), // The type of the interface
	}
}

var _ bson.ValueDecoder = (*TypedGenericCodex[string])(nil)
var _ bson.ValueEncoder = (*TypedGenericCodex[string])(nil)

// Register implements the Registrar interface.
// It allows the decoderEncoder to be registered with a bson.Registry.
func (t *TypedGenericCodex[Typed]) Register(builder Registrant) {

	builder.RegisterInterfaceEncoder(t.valueType, t)
	builder.RegisterInterfaceDecoder(t.valueType, t)
}

type wrapper struct {
	T string
	V bson.Raw
}

// RegisterType registers a factory function for a given type name
func (t *TypedGenericCodex[Typed]) RegisterType(
	factory func() Typed,
	toDto func(Typed) interface{},
	fromDto func(interface{}) Typed,
) {

	// check if the functions convert correctly
	value := factory()
	dto := toDto(value)
	converted := fromDto(dto)

	if !reflect.DeepEqual(value, converted) {
		panic(fmt.Errorf("toDto and fromDto functions do not convert correctly"))
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	type_ := t.getType(factory())
	t.subtypes[type_] = subtype[Typed]{factory: factory, toDto: toDto, fromDto: fromDto}
}

// lookupSubtype returns the factory function for a given type name
func (t *TypedGenericCodex[Typed]) lookupSubtype(typeName string) (subtype[Typed], bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	subtype, found := t.subtypes[typeName]
	return subtype, found
}

// EncodeValue implements the bson.ValueEncoder interface
func (t *TypedGenericCodex[Typed]) EncodeValue(ctx bson.EncodeContext, writer bson.ValueWriter, value reflect.Value) error {
	// Encode the original underlying value (it is the struct, not the interface)
	typed, ok := value.Interface().(Typed)
	if !ok {
		return fmt.Errorf("value does not implement Typed interface")
	}

	typeName := t.getType(typed)

	st, found := t.lookupSubtype(typeName)

	if !found {
		return fmt.Errorf("type %s not registered", typeName)
	}

	dto := st.toDto(typed)

	buf := new(bytes.Buffer)
	enc := bson.NewEncoder(bson.NewDocumentWriter(buf))
	enc.SetRegistry(ctx.Registry)
	enc.UseJSONStructTags()

	if err := enc.Encode(dto); err != nil {
		return err
	}

	// Wrap the original value with its type
	v := wrapper{T: typeName, V: buf.Bytes()}

	// Encode the wrapped value
	encoder, err := ctx.Registry.LookupEncoder(reflect.TypeOf(v))
	if err != nil {
		return err
	}

	return encoder.EncodeValue(ctx, writer, reflect.ValueOf(v))
}

// DecodeValue implements the bson.ValueDecoder interface
func (t *TypedGenericCodex[Typed]) DecodeValue(ctx bson.DecodeContext, reader bson.ValueReader, value reflect.Value) error {
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
	st, found := t.lookupSubtype(v.T)

	if !found {
		return fmt.Errorf("unknown type: %s", v.T)
	}

	// Decode the original underlying value
	result := st.factory()

	dto := st.toDto(result)

	dec := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(v.V)))
	dec.SetRegistry(ctx.Registry)
	dec.UseJSONStructTags()

	err = dec.Decode(dto)

	result = st.fromDto(dto)

	if err != nil {
		return err
	}

	// Set the original underlying value to the target interface
	value.Set(reflect.ValueOf(result))

	return nil
}
