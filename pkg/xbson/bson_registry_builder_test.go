package xbson

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Test_BuildDefaultRegistry_stores_registry_in_xbson_DefaultRegistry(t *testing.T) {
	saved := DefaultRegistry
	defer func() { DefaultRegistry = saved }()

	NewBsonRegistryBuilder().Build()

	assert.NotNil(t, DefaultRegistry, "Build() should populate xbson.DefaultRegistry")
}

func Test_BsonRegistryBuilder_register_codecs(t *testing.T) {
	type stringCodec struct{}
	enc := bson.ValueEncoderFunc(func(_ bson.EncodeContext, _ bson.ValueWriter, _ reflect.Value) error { return nil })
	dec := bson.ValueDecoderFunc(func(_ bson.DecodeContext, _ bson.ValueReader, _ reflect.Value) error { return nil })

	builder := NewBsonRegistryBuilder()
	builder.RegisterTypeEncoder(reflect.TypeOf(stringCodec{}), enc)
	builder.RegisterTypeDecoder(reflect.TypeOf(stringCodec{}), dec)

	got, err := builder.Registry().LookupEncoder(reflect.TypeOf(stringCodec{}))
	assert.NoError(t, err)
	assert.NotNil(t, got)

	gotDec, err := builder.Registry().LookupDecoder(reflect.TypeOf(stringCodec{}))
	assert.NoError(t, err)
	assert.NotNil(t, gotDec)
}
