package date

import (
	"time"
)

var (
	dateTimeZero      time.Time
	dateTimeFmt       = "2006-01-02T15:04:05.999999999Z"
	dateTimeFmtTicked = `"` + dateTimeFmt + `"`
)

// DateTime is timestamp in "date-time" format defined in RFC3339
type DateTime time.Time

// MarshalJSON override marshalJSON
func (dt DateTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(dt).Format(dateTimeFmtTicked)), nil
}

// UnmarshalJSON override unmarshalJSON
func (dt *DateTime) UnmarshalJSON(b []byte) error {
	ts, err := time.Parse(dateTimeFmtTicked, string(b))
	if err != nil {
		return err
	}

	*dt = DateTime(ts)
	return nil
}

// Reset set time zero value
func (dt *DateTime) Reset() {
	*dt = DateTime(dateTimeZero)
}

// String returns it's string representation
func (dt DateTime) String() string {
	return time.Time(dt).Format(dateTimeFmt)
}
