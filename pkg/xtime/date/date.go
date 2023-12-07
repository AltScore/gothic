package date

import (
	"fmt"
	"time"
)

const RFC3339Date = `2006-01-02`

// Empty is the zero value for a date. It is January 1, year 1, UTC.
var Empty = Date{}

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

// FromInLoc returns a new date from the given time in the given location. The time is truncated to the day.
func FromInLoc(t time.Time, loc *time.Location) Date {
	in := t.In(loc)
	return New(in.Year(), in.Month(), in.Day())
}

// FromIn returns a new date from the given time in the given location name. The time is truncated to the day.
func FromIn(t time.Time, locationName string) (Date, error) {
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		return Date{}, fmt.Errorf("invalid location %q: %w", locationName, err)
	}
	return FromInLoc(t, loc), nil
}

// FromInMust returns a new date from the given time in the given location name. The time is truncated to the day.
func FromInMust(t time.Time, locationName string) Date {
	d, err := FromIn(t, locationName)
	if err != nil {
		panic(err)
	}
	return d
}

// Today returns the current date in UTC. It is a convenience function for From(time.Now()).
// It is equivalent to the current day at 00:00 time.
func Today() Date {
	return From(time.Now())
}

// TodayInLoc returns the current date in the given location. It is a convenience function for From(time.Now().In(loc)).
func TodayInLoc(loc *time.Location) Date {
	return FromInLoc(time.Now(), loc)
}

// TodayIn returns the current date in the given location name. If locationName is invalid, an error is returned.
func TodayIn(locationName string) (Date, error) {
	return FromIn(time.Now(), locationName)
}

// MustTodayIn returns the current date in the given location name. If locationName is invalid, it panics.
func MustTodayIn(locationName string) Date {
	d, err := TodayIn(locationName)
	if err != nil {
		panic(err)
	}
	return d
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

// Sub returns the number of days between this date and the other date.
func (d Date) Sub(other Date) int {
	return int(d.Time().Sub(other.Time()).Hours() / 24)
}

// IsZero returns true if the date is the zero value.  January 1, year 1, UTC
func (d Date) IsZero() bool {
	return d.Time().IsZero()
}

// Format returns a textual representation of this date in the provided format.
func (d Date) Format(layout string) string {
	return d.Time().Format(layout)
}

// Weekday returns the day of the week specified by this date. Sunday is day 0.
func (d Date) Weekday() time.Weekday {
	return d.Time().Weekday()
}

// After reports whether the date is after the other.
func (d Date) After(other Date) bool {
	return d.IsAfter(other)
}

// Before reports whether this date is before the other.
func (d Date) Before(other Date) bool {
	return d.IsBefore(other)
}

// Equal reports whether this date is equal to the other.
// Deprecated: use Equals instead.
func (d Date) Equal(other Date) bool {
	return d.IsEqual(other)
}

// Equals reports whether this date is equal to the other.
func (d Date) Equals(other Date) bool {
	return d.IsEqual(other)
}

// NonZeroMin returns the earlier of the two dates
// ignoring zero values.
// If both are zero, returns Zero.
func (d Date) NonZeroMin(other Date) Date {
	if d.IsZero() {
		return other
	}
	if other.IsZero() {
		return d
	}
	return d.Min(other)
}

// Min returns the earlier of the two dates.
func (d Date) Min(other Date) Date {
	if d.Before(other) {
		return d
	}
	return other
}

// Max returns the latest of the two dates.
func (d Date) Max(other Date) Date {
	if d.After(other) {
		return d
	}
	return other
}

// Earliest returns the earliest of the dates.
func (d Date) Earliest(other Date) Date {
	return d.Min(other)
}

// Latest returns the latest of the dates.
func (d Date) Latest(other Date) Date {
	return d.Max(other)
}

// AsNullable returns the date as a nullable date (pointer to).
func (d Date) AsNullable() *Date {
	if d.IsZero() {
		return nil
	}
	return &d
}

// IsBefore reports whether the date is before the other date.
func (d Date) IsBefore(date Date) bool {
	return d.Time().Before(date.Time())
}

// IsAfter reports whether the date is after the other date.
func (d Date) IsAfter(date Date) bool {
	return d.Time().After(date.Time())
}

// IsEqual reports whether the date is equal to the other date.
func (d Date) IsEqual(date Date) bool {
	return d.Time().Equal(date.Time())
}

// IsNotEqual reports whether this date is not equal to the other.
func (d Date) IsNotEqual(other Date) bool {
	return !d.IsEqual(other)
}

// IsBetween reports whether the date is between the start and end dates (not inclusive).
func (d Date) IsBetween(start, end Date) bool {
	return d.IsAfter(start) && d.IsBefore(end)
}

// GoString returns a string representation of the date in Go syntax.
func (d Date) GoString() string {
	return fmt.Sprintf("date.Date(%d, %d, %d)", d.Year(), d.Month(), d.Day())
}

var _ fmt.GoStringer = Date{}
