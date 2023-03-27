package xerrors

import "testing"

func TestEnsureNotEmpty(t *testing.T) {
	emptyPointer := (*string)(nil)
	notEmptyPointer := new(string)

	tests := []struct {
		name string
		args any
		want bool
	}{
		{
			name: "empty string",
			args: "",
			want: true,
		},
		{
			name: "empty number",
			args: 0,
			want: true,
		},
		{
			name: "empty slice",
			args: []string{},
			want: true,
		},
		{
			name: "empty map",
			args: map[string]string{},
			want: true,
		},
		{
			name: "empty struct",
			args: struct{}{},
			want: false,
		},
		{
			name: "nil",

			args: nil,
			want: true,
		},
		{
			name: "empty pointer",

			args: emptyPointer,
			want: true,
		},
		{
			name: "non-empty string",
			args: "foo",
			want: false,
		},
		{
			name: "non-empty number",
			args: 42,
			want: false,
		},
		{
			name: "non-empty slice",
			args: []string{"foo"},
			want: false,
		},
		{
			name: "non-empty map",
			args: map[string]string{"foo": "bar"},
			want: false,
		},
		{
			name: "non-empty pointer",
			args: notEmptyPointer,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := runAndCatchPanic(func() {
				EnsureNotEmpty(tt.args, tt.name)
			})

			if actual != tt.want {
				t.Errorf("EnsureNotEmpty() = %v, want %v", actual, tt.want)
			}
		})
	}
}

func runAndCatchPanic(f func()) (caught bool) {
	defer func() {
		caught = recover() != nil
	}()
	f()
	return
}
