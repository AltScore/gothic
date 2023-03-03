package xvalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type sample struct {
	Amount string `validate:"required,numeric,vgt=0"`
}

func TestName(t *testing.T) {
	s := sample{
		Amount: "0",
	}

	err := Instance().Struct(s)

	assert.ErrorContains(t, err, "Key: 'sample.Amount' Error:Field validation for 'Amount' failed on the 'vgt' tag")
}

func TestStaticInstance(t *testing.T) {
	s := sample{
		Amount: "0",
	}

	err := Struct(s)

	assert.ErrorContains(t, err, "Key: 'sample.Amount' Error:Field validation for 'Amount' failed on the 'vgt' tag")
}

type sampleJSONArray struct {
	JsonArray string `validate:"json_array"`
}

func Test_isJSONArray(t *testing.T) {
	const errorOnJSONArrayTag = "Field validation for 'JsonArray' failed on the 'json_array' tag"

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "empty",
			value: "",
			want:  errorOnJSONArrayTag,
		},
		{
			name:  "spaces",
			value: " ",
			want:  errorOnJSONArrayTag,
		},
		{
			name:  "invalid",
			value: "invalid",
			want:  errorOnJSONArrayTag,
		},
		{
			name:  "invalid not array",
			value: `{"a": "b"}`,
			want:  errorOnJSONArrayTag,
		},
		{
			name:  "valid",
			value: "[\"a\", \"b\"]",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sampleJSONArray{
				JsonArray: tt.value,
			}

			err := Instance().Struct(s)

			if tt.want == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.want)
			}
		})
	}
}
