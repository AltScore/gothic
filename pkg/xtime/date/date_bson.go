package date

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var ErrInvalidDate = errors.New("invalid date")

func (d *Date) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if d == nil {
		return bson.ErrDecodeToNil
	}

	var err error

	switch t {
	case bsontype.Null:
		d.t = time.Time{}

	case bsontype.DateTime:
		if tm, _, ok := bsoncore.ReadTime(data); ok {
			d.t = From(tm).Time()
		} else {
			err = ErrInvalidDate
		}
	case bsontype.String:
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

func (d Date) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(d.Time())
}
