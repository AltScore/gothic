package ids

import (
	"reflect"
	"testing"
)

const (
	equal       = 0
	lowerThan   = -1
	greaterThan = 1
)

func TestCompare(t *testing.T) {
	type args struct {
		id    Id
		other Id
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "compare empty ids",
			args: args{
				id:    Empty(),
				other: Empty(),
			},
			want: equal,
		},
		{
			name: "compare empty id with non empty id",
			args: args{
				id:    Empty(),
				other: New(),
			},
			want: lowerThan,
		},
		{
			name: "compare non empty id with empty id",
			args: args{
				id:    New(),
				other: Empty(),
			},
			want: greaterThan,
		},
		{
			name: "compare non empty non equal ids",
			args: args{
				id:    MustParse("eb3401c8-190c-409b-b352-1a9529e1ca34"),
				other: MustParse("17d3ce07-0960-42a5-9b89-1f4a1638c9c5"),
			},
			want: greaterThan,
		},
		{
			name: "compare non empty equal ids",
			args: args{
				id:    MustParse("eb3401c8-190c-409b-b352-1a9529e1ca34"),
				other: MustParse("eb3401c8-190c-409b-b352-1a9529e1ca34"),
			},
			want: equal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compare(tt.args.id, tt.args.other); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrDefault(t *testing.T) {
	type args struct {
		id        Id
		defaultId Id
	}
	tests := []struct {
		name string
		args args
		want Id
	}{
		{
			name: "default id is returned when id is empty",
			args: args{
				id:        Empty(),
				defaultId: MustParse("7fc154a2-84a1-4a9f-a387-48de083a565e"),
			},
			want: MustParse("7fc154a2-84a1-4a9f-a387-48de083a565e"),
		},
		{
			name: "id is returned when id is not empty",
			args: args{
				id:        MustParse("d7eaba62-fd6e-412c-8707-9a71b9c5e7ec"),
				defaultId: MustParse("44ecbe0f-73e0-422b-a590-69985900dff6"),
			},
			want: MustParse("d7eaba62-fd6e-412c-8707-9a71b9c5e7ec"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrDefault(tt.args.id, tt.args.defaultId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
