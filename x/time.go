package x

import "time"

func ToNullableTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
