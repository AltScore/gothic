package date

import (
	"reflect"
	"testing"
	"time"
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
