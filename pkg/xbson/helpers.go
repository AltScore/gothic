package xbson

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

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
