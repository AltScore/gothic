package xmail

import "strings"

// Email represents an email address.
type Email string

// New creates a new Email.
func New(email string) Email {
	return Email(email)
}

func (e Email) IsValid() bool {
	return Validate(string(e))
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Trim() Email {
	return Email(strings.TrimSpace(e.String()))
}
