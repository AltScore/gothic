package xmongo

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

type EncoderDecoder interface {
	bsoncodec.ValueDecoder
	bsoncodec.ValueEncoder
	Register(builder *bsoncodec.RegistryBuilder)
}

var _ EncoderDecoder = (*decoderEncoder[int, int, int])(nil)

type decoderEncoder[Entity, Dto, Base any] struct {
	toDto   func(Entity) Dto
	fromDto func(Dto) Entity
}

func NewDecoderEncoder[Entity, Dto, Base any](toDto func(Entity) Dto, fromDto func(Dto) Entity) EncoderDecoder {
	return &decoderEncoder[Entity, Dto, Base]{toDto: toDto, fromDto: fromDto}
}

func (d *decoderEncoder[Entity, Dto, Base]) Register(builder *bsoncodec.RegistryBuilder) {
	entityType := reflect.TypeOf((*Entity)(nil)).Elem()

	baseType := reflect.TypeOf((*Base)(nil)).Elem()

	builder.RegisterTypeDecoder(entityType, d)
	builder.RegisterTypeEncoder(baseType, d)
}

func (d *decoderEncoder[Entity, Dto, Base]) DecodeValue(ctx bsoncodec.DecodeContext, reader bsonrw.ValueReader, value reflect.Value) error {
	var dto Dto
	decoder, err := ctx.Registry.LookupDecoder(reflect.TypeOf(dto))
	if err != nil {
		return err
	}

	if err := decoder.DecodeValue(ctx, reader, reflect.ValueOf(&dto).Elem()); err != nil {
		return err
	}

	value.Set(reflect.ValueOf(d.fromDto(dto)))
	return nil
}

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
