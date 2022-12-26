package date

import "time"

const RFC3339Date = `2006-01-02`

type Date struct { // as struct until we register the decoder / encoder, if not it uses time.Time decoder
	t time.Time
}

func New(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func From(t time.Time) Date {
	return Date{t.Truncate(24 * time.Hour).UTC()}
}

func Parse(s string) (Date, bool) {
	t, err := time.Parse(RFC3339Date, s)
	if err != nil {
		return Date{}, false
	}

	return From(t), true
}

func (d Date) Time() time.Time {
	return d.t
}

func (d Date) Year() int {
	return d.Time().Year()
}

func (d Date) Month() time.Month {
	return d.Time().Month()
}

func (d Date) Day() int {
	return d.Time().Day()
}

func (d Date) AddDate(years int, months int, days int) Date {
	return From(d.Time().AddDate(years, months, days))
}

func (d Date) AddDays(days int) Date {
	return d.AddDate(0, 0, days)
}

func (d Date) Add(duration time.Duration) Date {
	return From(d.Time().Add(duration))
}

func (d Date) IsZero() bool {
	return d.Time().IsZero()
}

func (d Date) Format(layout string) string {
	return d.Time().Format(layout)
}

func (d Date) Weekday() time.Weekday {
	return d.Time().Weekday()
}
