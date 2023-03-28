package date

import (
	"github.com/stretchr/testify/assert"
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
	// check that d is Date{} if and only if d.IsZero is true
	d := Date{}
	assert.Equal(t, d.IsZero(), true)
	d = From(time.Time{})
	assert.Equal(t, d.IsZero(), true)
	assert.Equal(t, d, Date{})
}

func Test_ZeroDateIsNotUnixZero(t *testing.T) {
	d := From(time.Unix(0, 0))
	assert.NotEqual(t, d, Date{})
	assert.Equal(t, d.IsZero(), false)
}

func TestDate_Min(t *testing.T) {
	d1 := Date{}
	d2 := Date{}
	assert.Equal(t, d1.Min(d2), Date{})

	d1 = Date{}
	d2 = New(2000, 10, 1)
	assert.Equal(t, d1.Min(d2), d1)
	assert.Equal(t, d2.Min(d1), d1)

	d1 = New(2001, 10, 1)
	d2 = New(2000, 10, 1)
	assert.Equal(t, d1.Min(d2), d2)
}

func TestDate_NonZeroMin(t *testing.T) {
	d1 := Date{}
	d2 := Date{}
	assert.Equal(t, d1.NonZeroMin(d2), Date{})

	d1 = Date{}
	d2 = New(2000, 10, 1)
	assert.Equal(t, d1.NonZeroMin(d2), d2)
	assert.Equal(t, d2.NonZeroMin(d1), d2)

	d1 = New(2001, 10, 1)
	d2 = New(2000, 10, 1)
	assert.Equal(t, d1.NonZeroMin(d2), d2)
}
