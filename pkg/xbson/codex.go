package xbson

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

type Registrar interface {
	Register(builder *bsoncodec.RegistryBuilder)
}

// EncoderDecoder is a bsoncodec.ValueDecoder and bsoncodec.ValueEncoder for a given type
type EncoderDecoder interface {
	Registrar
	bsoncodec.ValueDecoder
	bsoncodec.ValueEncoder
}

var _ EncoderDecoder = (*decoderEncoder[int, int, int])(nil)

// decoderEncoder is a bsoncodec.ValueDecoder and bsoncodec.ValueEncoder for a given type
// Entity is the interface type of the entity
// Dto is the type of the Data Transfer Object to store the entity
// Base is the type of the base entity that implements the Entity interface
type decoderEncoder[Entity, Dto, Base any] struct {
	toDto   func(Entity) Dto
	fromDto func(Dto) Entity
}

// NewDecoderEncoder creates a new decoderEncoder backed by the provided conversion functions
func NewDecoderEncoder[Entity, Dto, Base any](toDto func(Entity) Dto, fromDto func(Dto) Entity) EncoderDecoder {
	return &decoderEncoder[Entity, Dto, Base]{toDto: toDto, fromDto: fromDto}
}

// Register implements the bsoncodec.RegistryBuilder interface
// It allows the decoderEncoder to be registered with a bsoncodec.RegistryBuilder
func (d *decoderEncoder[Entity, Dto, Base]) Register(builder *bsoncodec.RegistryBuilder) {
	entityType := reflect.TypeOf((*Entity)(nil)).Elem()

	baseType := reflect.TypeOf((*Base)(nil)).Elem()

	builder.RegisterTypeEncoder(baseType, d)

	builder.RegisterHookDecoder(entityType, d)
	// builder.RegisterTypeDecoder(baseType, d)
}

// EncodeValue implements the bsoncodec.ValueEncoder interface. It encodes a Go value into a bson value
func (d *decoderEncoder[Entity, Dto, Base]) EncodeValue(ctx bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
	entity := value.Interface().(Entity)
	dto := d.toDto(entity)

	valueOfDto := reflect.ValueOf(dto)

	encoder, err := ctx.Registry.LookupEncoder(valueOfDto.Type())
	if err != nil {
		return err
	}

	return encoder.EncodeValue(ctx, writer, valueOfDto)
}

// DecodeValue implements the bsoncodec.ValueDecoder interface. It decodes a bson value into a Go value
func (d *decoderEncoder[Entity, Dto, Base]) DecodeValue(ctx bsoncodec.DecodeContext, reader bsonrw.ValueReader, value reflect.Value) error {
	var dto Dto
	decoder, err := ctx.Registry.LookupDecoder(reflect.TypeOf(dto))
	if err != nil {
		return err
	}

	if err := decoder.DecodeValue(ctx, reader, reflect.ValueOf(&dto).Elem()); err != nil {
		return err
	}

	fromDto := d.fromDto(dto)

	elem := reflect.ValueOf(fromDto).Elem()

	if elem.Type() == value.Type() {
		value.Set(elem)
	} else if elem.Type().Kind() != reflect.Ptr {
		value.Set(elem.Addr())
	} else {
		// value of type principal.principal is not assignable to type principal.Principal
		value.Set(elem)
	}

	return nil
}
