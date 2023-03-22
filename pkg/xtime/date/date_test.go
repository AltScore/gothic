package date

import (
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
