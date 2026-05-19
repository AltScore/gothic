package uuidbson

import (
	"fmt"
	"reflect"

	"github.com/AltScore/gothic/v2/pkg/xbson"
	guuid "github.com/google/uuid"
	"github.com/looplab/eventhorizon/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UUIDCodec2 struct {
	typeOfUUID reflect.Type
}

var _ xbson.EncoderDecoder = &UUIDCodec2{}

func (u *UUIDCodec2) Register(builder xbson.Registrant) {
	uuidType := reflect.TypeOf(uuid.Nil)

	u.typeOfUUID = uuidType

	builder.RegisterTypeEncoder(uuidType, u)
	builder.RegisterTypeDecoder(uuidType, u)
}

// EncodeValue Implement the ValueEncoder interface method.
func (u *UUIDCodec2) EncodeValue(_ bson.EncodeContext, vw bson.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != u.typeOfUUID || val.Len() != 16 {
		return bson.ValueEncoderError{
			Name:     "uuid.UUID",
			Types:    []reflect.Type{u.typeOfUUID},
			Received: val,
		}
	}
	b := make([]byte, 16)
	v := reflect.ValueOf(b)
	reflect.Copy(v, val)
	id, err := guuid.FromBytes(v.Bytes())
	if err != nil {
		return fmt.Errorf("could not parse UUID bytes (%x): %w", v.Bytes(), err)
	}

	return vw.WriteString(id.String())
}

// DecodeValue Implement the ValueDecoder interface method.
//
// Accepts both string and binary (subtype 4) BSON values for backward
// compatibility with documents written by mongo-driver v1, which encoded UUIDs
// as binary by default.
func (u *UUIDCodec2) DecodeValue(_ bson.DecodeContext, vr bson.ValueReader, val reflect.Value) error {
	if !val.IsValid() || !val.CanSet() || val.Kind() != reflect.Array {
		return bson.ValueDecoderError{
			Name:     "uuid.UUID",
			Kinds:    []reflect.Kind{reflect.Array},
			Received: val,
		}
	}

	var id uuid.UUID
	switch vr.Type() {
	case bson.TypeString:
		s, err := vr.ReadString()
		if err != nil {
			return err
		}

		id, err = uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("could not parse UUID string: %s", s)
		}
	case bson.TypeBinary:
		data, subtype, err := vr.ReadBinary()
		if err != nil {
			return err
		}

		if subtype != bson.TypeBinaryUUID {
			return fmt.Errorf("expected binary subtype 0x04 for UUID, got 0x%02x", subtype)
		}

		gid, err := guuid.FromBytes(data)
		if err != nil {
			return fmt.Errorf("could not parse UUID binary (%x): %w", data, err)
		}
		id = uuid.UUID(gid)
	default:
		return fmt.Errorf("received invalid BSON type to decode into UUID: %s", vr.Type())
	}

	v := reflect.ValueOf(id)
	if !v.IsValid() || v.Kind() != reflect.Array {
		return fmt.Errorf("invalid kind of reflected UUID value: %s", v.Kind().String())
	}
	reflect.Copy(val, v)

	return nil
}
