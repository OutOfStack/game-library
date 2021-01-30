package types

import (
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
