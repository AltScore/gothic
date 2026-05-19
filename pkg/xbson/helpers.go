package xbson

import (
	"bytes"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func UnmarshalWithRegistry(registry *bson.Registry, bs []byte, value interface{}) error {
	dec := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(bs)))
	dec.SetRegistry(registry)
	dec.UseJSONStructTags()

	return dec.Decode(value)
}

func MarshalWithRegistry(registry *bson.Registry, value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := bson.NewEncoder(bson.NewDocumentWriter(buf))
	enc.SetRegistry(registry)
	enc.UseJSONStructTags()

	if err := enc.Encode(value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
