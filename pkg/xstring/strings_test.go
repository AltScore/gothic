package xstring

import "testing"

func TestIsAllEmpty(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "empty string",
			args: "",
			want: true,
		},
		{
			name: "string with spaces",
			args: " ",
			want: true,
		},
		{
			name: "string with tabs",
			args: "\t",
			want: true,
		},
		{
			name: "string with newlines",
			args: "\n",
			want: true,
		},
		{
			name: "string with spaces and tabs",
			args: " \t",
			want: true,
		},
		{
			name: "string with spaces and newlines",
			args: " \n",
			want: true,
		},
		{
			name: "string with tabs and newlines",
			args: "\t\n",
			want: true,
		},
		{
			name: "string with spaces, tabs and newlines",
			args: " \t \n",
			want: true,
		},
		{
			name: "string with other characters",
			args: "a",
			want: false,
		},
		{
			name: "string with other characters and spaces",
			args: " a    ",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllEmpty(tt.args); got != tt.want {
				t.Errorf("IsAllEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsProperName(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "empty string",
			args: "",
			want: false,
		},
		{
			name: "string with spaces as prefix",
			args: " A name",
			want: false,
		},
		{
			name: "string with spaces as suffix",
			args: "A name ",
			want: false,
		},
		{
			name: "string with tabs",
			args: "A\tname",
			want: false,
		},
		{
			name: "string with newlines",
			args: "A\nnice name",
			want: false,
		},
		{
			name: "string with spaces, tabs and newlines",
			args: " A\tnice \nname ",
			want: false,
		},
		{
			name: "string proper name",
			args: "A nice name",
			want: true,
		},
		{
			name: "string proper name without spaces",
			args: "Anicename",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsProperName(tt.args); got != tt.want {
				t.Errorf("IsProperName() = %v, want %v", got, tt.want)
			}
		})
	}
}
