package xmail

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_atext_validate_correct_local_part(t *testing.T) {
	re := regexp.MustCompile(`^` + atext + `$`)

	valid := []string{
		"aaaaa",
		"bbb",
		"a-b",
		"a+b",
		"a1231239218",
	}

	for _, v := range valid {
		assert.True(t, re.MatchString(v))
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "valid email",
			email: "simple@domain.com",
			want:  true,
		},
		{
			name:  "valid email with subdomain",
			email: "another+sub@other.place.ai",
			want:  true,
		},
		{
			name:  "invalid email: no @",
			email: "simpledomain.com",
			want:  false,
		},
		{
			name:  "invalid email: no domain",
			email: "simple@",
			want:  false,
		},
		{
			name:  "valid email used in automatic UAT",
			email: "john.doe+UAT-a2-5689688D@altscore.ai",
			want:  true,
		},
		{
			name:  "invalid email used in automatic UAT (colon is not allowed in this format)",
			email: "john.doe+UAT:a2-5689688D@altscore.ai",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.email); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateList(t *testing.T) {

	tests := []struct {
		name   string
		emails string
		want   bool
	}{
		{
			name:   "valid list of one email",
			emails: "one@place.com.ar",
			want:   true,
		},
		{
			name:   "valid list of two emails",
			emails: "first@a.place.com,second@another.place.com",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateList(tt.emails); got != tt.want {
				t.Errorf("ValidateList() = %v, want %v", got, tt.want)
			}
		})
	}
}
