package date

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/x/bsonx/bsoncore"
)

var ErrInvalidDate = errors.New("invalid date")

func (d *Date) UnmarshalBSONValue(t byte, data []byte) error {
	if d == nil {
		return bson.ErrDecodeToNil
	}

	var err error

	switch bson.Type(t) {
	case bson.TypeNull:
		d.t = time.Time{}

	case bson.TypeDateTime:
		if tm, _, ok := bsoncore.ReadTime(data); ok {
			d.t = From(tm).Time()
		} else {
			err = ErrInvalidDate
		}
	case bson.TypeString:
		if date, ok := Parse(string(data)); ok {
			d.t = date.Time()
		} else {
			err = ErrInvalidDate
		}
	default:
		err = ErrInvalidDate
	}

	return err
}

func (d Date) MarshalBSONValue() (byte, []byte, error) {
	t, b, err := bson.MarshalValue(d.Time())
	return byte(t), b, err
}
