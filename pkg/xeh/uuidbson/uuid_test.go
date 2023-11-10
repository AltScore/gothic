package uuidbson

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/AltScore/gothic/v2/pkg/xbson"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUuid(t *testing.T) {
	base64String := "oR+7x682RYCPI+FCQYb5mA=="

	decodedBytes, err := base64.StdEncoding.DecodeString(base64String)

	require.NoError(t, err)

	anUuid, err := uuid.FromBytes(decodedBytes)

	require.NoError(t, err)

	fmt.Println(anUuid.String())
}

type sampleStructWithUuid struct {
	ID uuid.UUID
}

func TestCodec_encode_uuid_into_string_and_back(t *testing.T) {
	// GIVEN a UUID
	uuidStr := "b4e57d73-34ce-44b2-a57d-7334cea4b2d5"

	original := &sampleStructWithUuid{
		ID: uuid.MustParse(uuidStr),
	}

	// AND a bson registry
	registry := bson.NewRegistry()

	(&UUIDCodec2{}).Register(registry)

	// WHEN encoding it into bson
	bsonBytes, err := xbson.MarshalWithRegistry(registry, original)

	require.NoError(t, err)

	xbson.Dump(t, "UUDI", bsonBytes)

	// THEN the bson bytes should be a string
	var bsonMap bson.M

	err = bson.Unmarshal(bsonBytes, &bsonMap)

	require.NoError(t, err)

	require.Equal(t, uuidStr, bsonMap["id"])

	// WHEN decoding it back
	var decoded *sampleStructWithUuid

	err = xbson.UnmarshalWithRegistry(registry, bsonBytes, &decoded)

	require.NoError(t, err)

	// THEN the decoded object should be the same as the original
	require.Equal(t, original, decoded)
}

type sampleStructWithUuidPointer struct {
	ID *uuid.UUID
}

func TestCodec_encode_uuid_pointer_into_string_and_back(t *testing.T) {
	// GIVEN a UUID
	uuidStr := "b4e57d73-34ce-44b2-a57d-7334cea4b2d5"

	givenUuid := uuid.MustParse(uuidStr)

	original := &sampleStructWithUuidPointer{
		ID: &givenUuid,
	}

	// AND a bson registry
	registry := bson.NewRegistry()

	(&UUIDCodec2{}).Register(registry)

	// WHEN encoding it into bson
	bsonBytes, err := xbson.MarshalWithRegistry(registry, original)

	require.NoError(t, err)

	xbson.Dump(t, "UUDI", bsonBytes)

	// THEN the bson bytes should be a string
	var bsonMap bson.M

	err = bson.Unmarshal(bsonBytes, &bsonMap)

	require.NoError(t, err)

	require.Equal(t, uuidStr, bsonMap["id"])

	// WHEN decoding it back
	var decoded *sampleStructWithUuidPointer

	err = xbson.UnmarshalWithRegistry(registry, bsonBytes, &decoded)

	require.NoError(t, err)

	// THEN the decoded object should be the same as the original
	require.Equal(t, original, decoded)
}
