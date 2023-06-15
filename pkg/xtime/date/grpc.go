package date

import (
	"time"
)

// ProtoDate is a date proto interface to read Proto buffer google dates.
type ProtoDate interface {
	GetYear() int32
	GetMonth() int32
	GetDay() int32
}

// FromProto converts a google proto date to a date.
func FromProto(p ProtoDate) Date {
	if p == nil {
		return Date{}
	}
	return New(int(p.GetYear()), time.Month(p.GetMonth()), int(p.GetDay()))
}
