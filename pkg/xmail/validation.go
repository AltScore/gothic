package xmail

import "regexp"

const (
	atext = "[0-9a-zA-Z!#$%&'*+/=?^_`{|}~-]+"
	atom  = atext + `(\.` + atext + `)*`

	domain                     = `[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*`
	BaseFormatOfStringForEmail = atom + "@" + domain

	FormatOfStringForEmail     = `^` + BaseFormatOfStringForEmail + `$`
	FormatOfStringForEmailList = `^` + BaseFormatOfStringForEmail + `([,;]` + BaseFormatOfStringForEmail + `)*$`
)

var (
	emailRe     = regexp.MustCompile(FormatOfStringForEmail)
	emailListRe = regexp.MustCompile(FormatOfStringForEmailList)
)

// Validate checks if the given string is a valid email address.
func Validate(email string) bool {
	return emailRe.MatchString(email)
}

// ValidateList checks if the given string is a valid list of email addresses.
func ValidateList(emails string) bool {
	return emailListRe.MatchString(emails)
}
