package xbson

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

// BsonRegistryBuilder initializes the mongo driver registry to encode/decode honoring the JSON struct tags.
// For example, a struct with the following JSON tags:
//
//	type Sample struct {
//	       FirstName string `json:"first_name"`
//	}
//
// Will serialize the field to BSON as "first_name" instead of "firstname" (default naming strategy).
type BsonRegistryBuilder struct {
	*bsoncodec.RegistryBuilder
	structCodec *bsoncodec.StructCodec
}

type BsonCodecsRegistrant func(builder *BsonRegistryBuilder)

var DefaultBsonRegistryBuilder = NewBsonRegistryBuilder()

func NewBsonRegistryBuilder() *BsonRegistryBuilder {
	codec, err := bsoncodec.NewStructCodec(bsoncodec.JSONFallbackStructTagParser)

	if err != nil {
		panic(err)
	}

	builder := bson.NewRegistryBuilder()
	builder.RegisterDefaultEncoder(reflect.Struct, codec)
	builder.RegisterDefaultDecoder(reflect.Struct, codec)

	return &BsonRegistryBuilder{
		RegistryBuilder: builder,
		structCodec:     codec,
	}
}

// Register a custom codec to the default BSON registry
func (b *BsonRegistryBuilder) Register(registrant BsonCodecsRegistrant) *BsonRegistryBuilder {
	registrant(b)
	return b
}

// RegisterAll register all the custom codecs to the default BSON registry
func (b *BsonRegistryBuilder) RegisterAll(registrants ...BsonCodecsRegistrant) *BsonRegistryBuilder {
	for _, registrant := range registrants {
		b.Register(registrant)
	}
	return b
}

func (b *BsonRegistryBuilder) RegisterTypeDecoder(t reflect.Type, dec bsoncodec.ValueDecoder) {
	b.RegistryBuilder.RegisterTypeDecoder(t, dec)
}

func (b *BsonRegistryBuilder) RegisterTypeEncoder(t reflect.Type, dec bsoncodec.ValueEncoder) {
	b.RegistryBuilder.RegisterTypeEncoder(t, dec)
}

// Build sets this registry as the BSON default
func (b *BsonRegistryBuilder) Build() {
	bson.DefaultRegistry = b.RegistryBuilder.Build()
}

// StructCodec provides the configured bsoncodec.StructCodec in registry
func (b *BsonRegistryBuilder) StructCodec() *bsoncodec.StructCodec {
	return b.structCodec
}
