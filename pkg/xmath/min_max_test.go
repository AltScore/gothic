package xmath

import (
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestMinInt(t *testing.T) {
	type args[N constraints.Ordered] struct {
		a N
		b N
	}
	type testCase[N constraints.Ordered] struct {
		name string
		args args[N]
		want N
	}
	tests := []testCase[int]{
		{
			name: "first",
			args: args[int]{a: 1, b: 2},
			want: 1,
		},
		{
			name: "second",
			args: args[int]{a: 2, b: 1},
			want: 1,
		},
		{
			name: "equal",
			args: args[int]{a: 123, b: 123},
			want: 123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMinSeveralInts(t *testing.T) {
	type testCase[N constraints.Ordered] struct {
		name string
		args []N
		want N
	}
	tests := []testCase[int]{
		{
			name: "only one",
			args: []int{79},
			want: 79,
		},
		{
			name: "first",
			args: []int{1, 2, 3, 6, 5},
			want: 1,
		},
		{
			name: "second",
			args: []int{3, 1, 2, 3, 6, 5},
			want: 1,
		},
		{
			name: "equal",
			args: []int{123, 123, 123, 123, 123},
			want: 123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args[0], tt.args[1:]...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args[N constraints.Ordered] struct {
		a N
		b N
	}
	type testCase[N constraints.Ordered] struct {
		name string
		args args[N]
		want N
	}
	tests := []testCase[float64]{
		{
			name: "first",
			args: args[float64]{a: 223.45, b: 112.24},
			want: 223.45,
		},
		{
			name: "second",
			args: args[float64]{a: 112.24, b: 223.45},
			want: 223.45,
		},
		{
			name: "equal",
			args: args[float64]{a: 123.45, b: 123.45},
			want: 123.45,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	type args[N constraints.Ordered] struct {
		value N
		min   N
		max   N
	}
	type testCase[N constraints.Ordered] struct {
		name string
		args args[N]
		want N
	}
	tests := []testCase[uint]{
		{
			name: "value less than min",
			args: args[uint]{value: 19, min: 20, max: 30},
			want: 20,
		},
		{
			name: "value greater than max",
			args: args[uint]{value: 31, min: 20, max: 30},
			want: 30,
		},
		{
			name: "value far less than min",
			args: args[uint]{value: 3, min: 20, max: 30},
			want: 20,
		},
		{
			name: "value far greater than max",
			args: args[uint]{value: 223, min: 20, max: 30},
			want: 30,
		},
		{
			name: "value between min and max",
			args: args[uint]{value: 25, min: 20, max: 30},
			want: 25,
		},
		{
			name: "value equal to min",
			args: args[uint]{value: 20, min: 20, max: 30},
			want: 20,
		},
		{
			name: "value equal to max",
			args: args[uint]{value: 30, min: 20, max: 30},
			want: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.args.value, tt.args.min, tt.args.max); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}
