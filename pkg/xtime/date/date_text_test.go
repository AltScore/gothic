package date

import (
	"reflect"
	"testing"
	"time"
)

func TestDate_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		d       Date
		want    []byte
		wantErr bool
	}{
		{
			name: "MarshalJSON",
			d:    New(2019, 1, 1),
			want: []byte(`2019-01-01`),
		},
		{
			name: "MarshalJSON 2",
			d:    From(time.Date(1963, 11, 29, 12, 23, 45, 0, time.UTC)),
			want: []byte(`1963-11-29`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		want    Date
		args    []byte
		wantErr bool
	}{
		{
			name: "UnmarshalText",
			args: []byte(`2019-04-19`),
			want: New(2019, 4, 19),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date
			if err := d.UnmarshalText(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(d, tt.want) {
				t.Errorf("UnmarshalText() got = %v, want %v", d, tt.want)
			}
		})
	}
}
