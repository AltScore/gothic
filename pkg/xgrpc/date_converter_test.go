package xgrpc

import (
	xdate "github.com/AltScore/gothic/v2/pkg/xtime/date"
	"google.golang.org/genproto/googleapis/type/date"
	"reflect"
	"testing"
)

func TestDateToProto(t *testing.T) {
	tests := []struct {
		name string
		args xdate.Date
		want *date.Date
	}{
		{
			name: "non empty date",
			args: xdate.New(2019, 11, 29),
			want: &date.Date{Year: 2019, Month: 11, Day: 29},
		},
		{
			name: "empty date",
			args: xdate.Date{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateToProto(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}
