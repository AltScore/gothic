package date

const RFC3339DateJSON = `"` + RFC3339Date + `"`

// MarshalJSON implements the json.Marshaler interface.
// Converts to "yyyy-MMM-dd" format.
func (d Date) MarshalJSON() ([]byte, error) {
	return d.MarshalText()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Converts from "yyyy-MMM-dd" format.
func (d *Date) UnmarshalJSON(b []byte) error {
	return d.UnmarshalText(b)
}
