package xgrpc

import (
	xdate "github.com/AltScore/gothic/v2/pkg/xtime/date"
	"google.golang.org/genproto/googleapis/type/date"
)

// DateToProto converts a gothic date.Date to a google proto date.Date
func DateToProto(d xdate.Date) *date.Date {
	if d.IsZero() {
		return nil
	}

	return &date.Date{
		Year:  int32(d.Year()),
		Month: int32(d.Month()),
		Day:   int32(d.Day()),
	}
}
