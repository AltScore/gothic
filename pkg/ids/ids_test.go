package ids

import "testing"

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
