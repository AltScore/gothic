package date

import (
	"time"
)

const RFC3339Date = `2006-01-02`

// Date represents a date without time. It is a wrapper around time.Time.
type Date struct { // as struct until we register the decoder / encoder, if not it uses time.Time decoder
	t time.Time
}

// New returns a new date with the given year, month and day.
func New(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// From returns a new date from the given time. The time is truncated to the day.
func From(t time.Time) Date {
	return Date{t.Truncate(24 * time.Hour).UTC()}
}

// Today returns the current date in UTC. It is a convenience function for From(time.Now()).
// It is equivalent to the current day at 00:00 time.
func Today() Date {
	return From(time.Now())
}

// Parse parses a date in RFC3339 format (yyyy-mm-dd).
func Parse(s string) (Date, bool) {
	t, err := time.Parse(RFC3339Date, s)
	if err != nil {
		return Date{}, false
	}

	return From(t), true
}

// MustParse parses a date in RFC3339 format (yyyy-mm-dd) and panics if the date is invalid.
// It is intended for use in variable initializations and tests.
func MustParse(s string) Date {
	d, ok := Parse(s)
	if !ok {
		panic("invalid date " + s)
	}

	return d
}
func (d Date) String() string {
	return d.Format(RFC3339Date)
}

// Time returns the time.Time representation of the date.
func (d Date) Time() time.Time {
	return d.t
}

// Year returns the year of the date.
func (d Date) Year() int {
	return d.Time().Year()
}

// Month returns the month of the date.
func (d Date) Month() time.Month {
	return d.Time().Month()
}

// Day returns the day of the date.
func (d Date) Day() int {
	return d.Time().Day()
}

// AddDate returns the date with the given number of years, months and days added.
// To subtract, use negative values.
func (d Date) AddDate(years int, months int, days int) Date {
	return From(d.Time().AddDate(years, months, days))
}

// AddDays returns the date with the given number of days added.
// To subtract, use negative values.
func (d Date) AddDays(days int) Date {
	return d.AddDate(0, 0, days)
}

// Add returns the date with the given duration added.
// Only the date part is used, the time part is ignored.
func (d Date) Add(duration time.Duration) Date {
	return From(d.Time().Add(duration))
}

// IsZero returns true if the date is the zero value.  January 1, year 1, UTC
func (d Date) IsZero() bool {
	return d.Time().IsZero()
}

// Format returns a textual representation of the date in the provided format.
func (d Date) Format(layout string) string {
	return d.Time().Format(layout)
}

// Weekday returns the day of the week specified by the date. Sunday is day 0.
func (d Date) Weekday() time.Weekday {
	return d.Time().Weekday()
}

// After reports whether the date is after the other.
func (d Date) After(other Date) bool {
	return d.Time().After(other.Time())
}

// Before reports whether the date is before the other.
func (d Date) Before(other Date) bool {
	return d.Time().Before(other.Time())
}

// Equal reports whether the date is equal to the other.
func (d Date) Equal(other Date) bool {
	return d.Time().Equal(other.Time())
}

func (d Date) NonZeroMin(other Date) Date {
	if d.IsZero() {
		return other
	}
	if other.IsZero() {
		return d
	}
	return d.Min(other)
}

func (d Date) Min(other Date) Date {
	if d.Before(other) {
		return d
	}
	return other
}

func (d Date) AsNullable() *Date {
	if d.IsZero() {
		return nil
	}
	return &d
}
