package date

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		d       Date
		want    []byte
		wantErr bool
	}{
		{
			name: "MarshalJSON",
			d:    New(2019, 1, 1),
			want: []byte(`"2019-01-01"`),
		},
		{
			name: "MarshalJSON 2",
			d:    From(time.Date(1963, 11, 29, 12, 23, 45, 0, time.UTC)),
			want: []byte(`"1963-11-29"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.d)
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

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		want    Date
		args    []byte
		wantErr bool
	}{
		{
			name: "UnmarshalJSON",
			args: []byte(`"2019-01-01"`),
			want: New(2019, 1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date
			if err := json.Unmarshal(tt.args, &d); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(d, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", d, tt.want)
			}
		})
	}
}

func Test_MarshallTime(t *testing.T) {
	t.Skip("We need to fix this in some way")
	tests := []struct {
		name    string
		want    []byte
		args    time.Time
		wantErr bool
	}{
		{
			name: "1 ns",
			args: time.Date(2019, 1, 1, 0, 0, 0, 123100000, time.UTC),
			want: []byte(`"2019-01-01T00:00:00.123Z"`),
		},
		{
			name: "0.1 ms",
			args: time.Date(2019, 1, 1, 0, 0, 0, 123000001, time.UTC),
			want: []byte(`"2019-01-01T00:00:00.123Z"`),
		},
		{
			name: "10 ms",
			args: time.Date(2019, 1, 1, 0, 0, 0, 10000000, time.UTC),
			want: []byte(`"2019-01-01T00:00:00.123Z"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshallTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshallTime() got = %s, want %s", got, tt.want)
			}
		})
	}
}
