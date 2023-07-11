package date

import (
	"time"
)

// MarshalText Implements encoding.TextMarshaler for Date
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.Time().Format(RFC3339Date)), nil
}

// UnmarshalText Implements encoding.TextUnmarshaler for Date
func (d *Date) UnmarshalText(text []byte) error {
	t, err := time.Parse(RFC3339Date, string(text))
	if err != nil {
		return err
	}
	d.t = From(t).Time()
	return nil
}
