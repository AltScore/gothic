package xstring

import "strings"

// IsAllEmpty returns true if the string is empty or contains only whitespace.
func IsAllEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsProperName returns true if the string is not empty and contains no whitespace
// prefix or suffix and does not contain tabs or newlines.
func IsProperName(s string) bool {
	return s != "" &&
		strings.TrimSpace(s) == s &&
		!strings.ContainsAny(s, "\t\n")
}

// IsNotProperName returns true if the string is empty, contains whitespace as
// prefix or suffix or contains tabs or newlines.
func IsNotProperName(s string) bool {
	return !IsProperName(s)
}
