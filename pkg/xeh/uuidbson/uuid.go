package uuidbson

import (
	"reflect"

	"github.com/AltScore/gothic/v2/pkg/xbson"
	"github.com/looplab/eventhorizon/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

type UUIDCodec struct {
	typeOfUUID reflect.Type
}

var _ xbson.EncoderDecoder = &UUIDCodec{}

func (u *UUIDCodec) Register(builder xbson.Registrant) {
	uuidType := reflect.TypeOf(uuid.Nil)

	u.typeOfUUID = uuidType

	builder.RegisterTypeEncoder(uuidType, u)
	builder.RegisterTypeDecoder(uuidType, u)
}

// EncodeValue Implement the ValueEncoder interface method.
func (u *UUIDCodec) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.Type().AssignableTo(reflect.TypeOf(uuid.UUID{})) {
		return bsoncodec.ValueEncoderError{Name: "UUIDCodec.EncodeValue", Types: []reflect.Type{u.typeOfUUID}, Received: val}
	}

	uu, ok := val.Interface().(uuid.UUID)

	if !ok {
		return bsoncodec.ValueEncoderError{Name: "UUIDCodec.EncodeValue", Types: []reflect.Type{u.typeOfUUID}, Received: val}
	}

	return vw.WriteString(uu.String())
}

// DecodeValue Implement the ValueDecoder interface method.
func (u *UUIDCodec) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != u.typeOfUUID {
		return bsoncodec.ValueDecoderError{Name: "UUIDCodec.DecodeValue", Types: []reflect.Type{u.typeOfUUID}, Received: val}
	}

	if vr.Type() != bson.TypeString {
		return bsoncodec.ValueDecoderError{Name: "UUIDCodec.DecodeValue", Types: []reflect.Type{u.typeOfUUID}, Received: val}
	}

	str, err := vr.ReadString()
	if err != nil {
		return err
	}

	uu, err := uuid.Parse(str)
	if err != nil {
		return err
	}

	val.Set(reflect.ValueOf(uu))
	return nil
}
