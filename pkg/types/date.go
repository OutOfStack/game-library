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
	date := civil.DateOf(v.(time.Time))
	*t = Date(date)
	return nil
}

// String returns the date in RFC3339 full-date format.
func (t Date) String() string {
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
