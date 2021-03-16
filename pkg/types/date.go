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
