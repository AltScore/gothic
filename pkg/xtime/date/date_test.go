package date

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDate_After(t *testing.T) {
	tests := []struct {
		name   string
		fields Date
		args   Date
		want   bool
	}{
		{
			name:   "is before",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			want:   false,
		},
		{
			name:   "is after",
			fields: Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   true,
		},
		{
			name:   "is equal",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields
			if got := d.After(tt.args); got != tt.want {
				t.Errorf("After() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Before(t *testing.T) {

	tests := []struct {
		name   string
		fields Date
		args   Date
		want   bool
	}{
		{
			name:   "is before",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			want:   true,
		},
		{
			name:   "is after",
			fields: Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   false,
		},
		{
			name:   "is equal",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields
			if got := d.Before(tt.args); got != tt.want {
				t.Errorf("Before() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Equal(t *testing.T) {
	tests := []struct {
		name   string
		fields Date
		args   Date
		want   bool
	}{
		{
			name:   "is before",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			want:   false,
		},
		{
			name:   "is after",
			fields: Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   false,
		},
		{
			name:   "is equal",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields
			if got := d.Equal(tt.args); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsZero(t *testing.T) {
	tests := []struct {
		name   string
		fields Date
		want   bool
	}{
		{
			name:   "is zero",
			fields: Date{t: time.Time{}},
			want:   true,
		},
		{
			name:   "is empty",
			fields: Date{},
			want:   true,
		},
		{
			name:   "is not zero",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields
			if got := d.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestDate_Min(t *testing.T) {
	tests := []struct {
		name   string
		fields Date
		args   Date
		want   Date
	}{
		{
			name:   "is before",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			want:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:   "is after",
			fields: Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:   "is equal",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields
			if got := d.Min(tt.args); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_NonZeroMin(t *testing.T) {
	tests := []struct {
		name   string
		fields Date
		args   Date
		want   Date
	}{
		{
			name:   "is before",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			want:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:   "field is empty",
			fields: Date{},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:   "is after",
			fields: Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:   "arg is empty",
			fields: Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
			args:   Date{},
			want:   Date{t: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:   "is equal",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			args:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:   "field and args are empty",
			fields: Date{},
			args:   Date{},
			want:   Date{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields
			if got := d.NonZeroMin(tt.args); got != tt.want {
				t.Errorf("NonZeroMin() = %v, want %v", got, tt.want)
			}
		})
	}
}

// AsNullable() returns *Date or nil if the Date is empty.
func TestDate_AsNullable(t *testing.T) {
	tests := []struct {
		name   string
		fields Date
		want   *Date
	}{
		{
			name:   "is not empty",
			fields: Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:   &Date{t: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:   "is empty",
			fields: Date{},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields
			if got := d.AsNullable(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsNullable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_String(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want string
	}{
		{
			name: "is empty",
			d:    Date{},
			want: "0001-01-01",
		},
		{
			name: "is not empty",
			d:    New(2020, 8, 17),
			want: "2020-08-17",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Sub(t *testing.T) {
	tests := []struct {
		name string
		from Date
		args Date
		want int
	}{
		{
			name: "is before",
			from: New(2020, 1, 1),
			args: New(2020, 1, 2),
			want: -1,
		},
		{
			name: "is after",
			from: New(2020, 1, 2),
			args: New(2020, 1, 1),
			want: 1,
		},
		{
			name: "is equal",
			from: New(2020, 1, 1),
			args: New(2020, 1, 1),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.Sub(tt.args); got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FromProto_builds_date_with_correct_values(t *testing.T) {
	date := FromProto(testProto{})

	assert.Equal(t, 1963, date.Year())
	assert.Equal(t, time.November, date.Month())
	assert.Equal(t, 29, date.Day())
}

type testProto struct{}

var _ ProtoDate = (*testProto)(nil)

func (t testProto) GetYear() int32  { return 1963 }
func (t testProto) GetMonth() int32 { return 11 }

func (t testProto) GetDay() int32 { return 29 }

func TestFromInLoc(t *testing.T) {
	mexico := time.FixedZone("America/Mexico_City", -6*60*60)
	argentina := time.FixedZone("America/Argentina/Buenos_Aires", -3*60*60)
	chile := time.FixedZone("America/Santiago", -4*60*60)
	ecuador := time.FixedZone("America/Guayaquil", -5*60*60)

	type args struct {
		t   time.Time
		loc *time.Location
	}
	tests := []struct {
		name    string
		args    args
		want    Date
		wantStr string
	}{
		{
			name: "morning of Argentina in Mexico",
			args: args{
				t:   time.Date(2020, 7, 1, 2, 42, 28, 0, argentina),
				loc: mexico,
			},
			want:    New(2020, 6, 30),
			wantStr: "2020-06-30",
		},
		{
			name: "morning of Ecuador in Mexico",
			args: args{
				t:   time.Date(2020, 7, 1, 4, 42, 28, 0, ecuador),
				loc: mexico,
			},
			want:    New(2020, 7, 1),
			wantStr: "2020-07-01",
		},
		{
			name: "Morning of Ecuador in Chile",
			args: args{
				t:   time.Date(2020, 7, 1, 4, 42, 28, 0, ecuador),
				loc: chile,
			},
			want:    New(2020, 7, 1),
			wantStr: "2020-07-01",
		},
		{
			name: "Morning of UTC in Mexico",
			args: args{
				t:   time.Date(2020, 7, 1, 5, 59, 59, 0, time.UTC),
				loc: mexico,
			},
			want:    New(2020, 6, 30),
			wantStr: "2020-06-30",
		},
		{
			name: "Late Morning of UTC in Mexico",
			args: args{
				t:   time.Date(2020, 7, 1, 6, 0, 0, 0, time.UTC),
				loc: mexico,
			},
			want:    New(2020, 7, 1),
			wantStr: "2020-07-01",
		},
		{
			name: "Later Morning of UTC in Mexico",
			args: args{
				t:   time.Date(2020, 7, 1, 6, 0, 0, 1, time.UTC),
				loc: mexico,
			},
			want:    New(2020, 7, 1),
			wantStr: "2020-07-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := FromInLoc(tt.args.t, tt.args.loc)
			assert.Equalf(t, tt.want, d, "FromInLoc(%v, %v)", tt.args.t, tt.args.loc)
			assert.Equalf(t, tt.wantStr, d.String(), "FromInLoc(%v, %v)", tt.args.t, tt.args.loc)
		})
	}
}
