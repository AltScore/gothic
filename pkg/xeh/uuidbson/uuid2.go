package uuidbson

import (
	"fmt"
	"reflect"

	"github.com/AltScore/gothic/v2/pkg/xbson"
	guuid "github.com/google/uuid"
	"github.com/looplab/eventhorizon/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
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
func (u *UUIDCodec2) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != u.typeOfUUID || val.Len() != 16 {
		return bsoncodec.ValueEncoderError{
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
func (u *UUIDCodec2) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || !val.CanSet() || val.Kind() != reflect.Array {
		return bsoncodec.ValueDecoderError{
			Name:     "uuid.UUID",
			Kinds:    []reflect.Kind{reflect.Bool},
			Received: val,
		}
	}

	var s string
	switch vr.Type() {
	case bson.TypeString:
		var err error
		if s, err = vr.ReadString(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("received invalid BSON type to decode into UUID: %s", vr.Type())
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("could not parse UUID string: %s", s)
	}
	v := reflect.ValueOf(id)
	if !v.IsValid() || v.Kind() != reflect.Array {
		return fmt.Errorf("invalid kind of reflected UUID value: %s", v.Kind().String())
	}
	reflect.Copy(val, v)

	return nil
}
