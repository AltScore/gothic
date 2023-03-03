package xlogger

import "testing"

func TestHideKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty key",
			args: args{key: ""},
			want: "<NONE>",
		},
		{
			name: "one char key",
			args: args{key: "a"},
			want: "*",
		},
		{
			name: "two char key",
			args: args{key: "ab"},
			want: "ab...ab",
		},
		{
			name: "three char key",
			args: args{key: "abc"},
			want: "ab...bc",
		},
		{
			name: "four char key",
			args: args{key: "abcd"},
			want: "ab...cd",
		},
		{
			name: "five char key",
			args: args{key: "abcde"},
			want: "ab...de",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HideKey(tt.args.key); got != tt.want {
				t.Errorf("HideKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHideMongoSecret(t *testing.T) {
	type args struct {
		mongoUri string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty mongo uri",
			args: args{mongoUri: ""},
			want: "",
		},
		{
			name: "no password",
			args: args{mongoUri: "mongodb://localhost:27017"},
			want: "mongodb://localhost:27017",
		},
		{
			name: "password",
			args: args{mongoUri: "mongodb://user:password@localhost:27017"},
			want: "mongodb://user:***@localhost:27017",
		},
		{
			name: "password with special chars",
			args: args{mongoUri: "mongodb://user:pass@word@localhost:27017"},
			want: "mongodb://user:***@localhost:27017",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HideMongoSecret(tt.args.mongoUri); got != tt.want {
				t.Errorf("HideMongoSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
