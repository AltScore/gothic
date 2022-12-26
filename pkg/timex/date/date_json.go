package date

import "time"

const RFC3339DateJSON = `"` + RFC3339Date + `"`

// MarshalJSON implements the json.Marshaler interface.
// Converts to "yyyy-MMM-dd" format.
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Time().Format(RFC3339DateJSON)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Converts from "yyyy-MMM-dd" format.
func (d *Date) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(RFC3339DateJSON, string(b))
	if err != nil {
		return err
	}
	d.t = From(t).Time()
	return nil
}
