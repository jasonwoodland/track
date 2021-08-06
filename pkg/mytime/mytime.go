package mytime

import (
	"time"
)

type Time time.Time

// Allow scanning into date from using (*sql.Rows).Scan
func (t *Time) Scan(v interface{}) error {
	vt, err := time.Parse(time.RFC3339, string(v.(string)))
	if err != nil {
		return err
	}
	*t = Time(vt)
	return nil
}
