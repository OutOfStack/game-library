package types

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

// Date is date type
type Date civil.Date

// Scan casts time.Time date to Date type
func (t *Date) Scan(v interface{}) error {
	if v == nil {
		*t = Date(civil.Date{})
		return nil
	}
	d, ok := v.(time.Time)
	if !ok {
		return fmt.Errorf("parsing date to time.Time")
	}
	date := civil.DateOf(d)
	*t = Date(date)
	return nil
}

// String returns the date in RFC3339 full-date format.
func (t Date) String() string {
	d := civil.Date(t)
	if d.IsZero() {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", t.Year, t.Month, t.Day)
}

// ParseDate parses a string in 'YYYY-MM-DD' format and returns the date value it represents.
func ParseDate(s string) (Date, error) {
	d, err := civil.ParseDate(s)
	if err != nil {
		return Date{}, err
	}
	return Date(d), nil
}
