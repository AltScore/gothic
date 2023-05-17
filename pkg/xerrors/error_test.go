package xerrors

import (
	"testing"
)

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "not found",
			err:  NewNotFoundError("entity", "keyFmt %s", "args"),
			want: true,
		},
		{
			name: "other error",
			err:  NewInvalidArgumentError("entity", "keyFmt %s", "args"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}
