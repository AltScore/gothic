package xbson

import (
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BsonRegistryBuilder wraps a *bson.Registry to ease registration of custom codecs.
//
// Note: in mongo-driver v2 there is no longer a settable global default registry,
// nor a registry-level option to honor JSON struct tags. JSON-tag fallback is
// configured per encoder/decoder via *bson.Encoder.UseJSONStructTags() and
// *bson.Decoder.UseJSONStructTags(). Build() stores the registry in
// xbson.DefaultRegistry — consumers must pass it explicitly to their mongo
// client/database via options.Client().SetRegistry(...).
type BsonRegistryBuilder struct {
	registry *bson.Registry
}

type BsonCodecsRegistrant func(builder *BsonRegistryBuilder)

var DefaultBsonRegistryBuilder = NewBsonRegistryBuilder()

func NewBsonRegistryBuilder() *BsonRegistryBuilder {
	return &BsonRegistryBuilder{
		registry: bson.NewRegistry(),
	}
}

// Register a custom codec to the BSON registry
func (b *BsonRegistryBuilder) Register(registrant BsonCodecsRegistrant) *BsonRegistryBuilder {
	registrant(b)
	return b
}

// RegisterAll register all the custom codecs to the BSON registry
func (b *BsonRegistryBuilder) RegisterAll(registrants ...BsonCodecsRegistrant) *BsonRegistryBuilder {
	for _, registrant := range registrants {
		b.Register(registrant)
	}
	return b
}

func (b *BsonRegistryBuilder) RegisterTypeDecoder(t reflect.Type, dec bson.ValueDecoder) {
	b.registry.RegisterTypeDecoder(t, dec)
}

func (b *BsonRegistryBuilder) RegisterTypeEncoder(t reflect.Type, enc bson.ValueEncoder) {
	b.registry.RegisterTypeEncoder(t, enc)
}

func (b *BsonRegistryBuilder) RegisterInterfaceEncoder(t reflect.Type, enc bson.ValueEncoder) {
	b.registry.RegisterInterfaceEncoder(t, enc)
}

func (b *BsonRegistryBuilder) RegisterInterfaceDecoder(t reflect.Type, dec bson.ValueDecoder) {
	b.registry.RegisterInterfaceDecoder(t, dec)
}

// Build stores the registry in xbson.DefaultRegistry. Mongo-driver v2 has no
// settable global default registry, so consumers must pass xbson.DefaultRegistry
// explicitly via options.Client().SetRegistry(...) or
// options.Database().SetRegistry(...).
func (b *BsonRegistryBuilder) Build() {
	DefaultRegistry = b.registry
}

// Registry returns the underlying *bson.Registry for direct use.
func (b *BsonRegistryBuilder) Registry() *bson.Registry {
	return b.registry
}
