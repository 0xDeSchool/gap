package x

import (
	"encoding/json"
	"time"
)

// NullableTime is a wrapper of time.Time, which can be marshaled to JSON as RFC3339 format.
// It can be used as a nullable time.Time field.
// empty value is zero time
// value is null or time
type NullableTime struct {
	t        time.Time
	nilValue bool
}

func (t NullableTime) MarshalJSON() ([]byte, error) {
	if t.nilValue {
		return []byte("null"), nil
	} else {
		return json.Marshal(t.t)
	}
}

func (t *NullableTime) UnmarshalJSON(data []byte) error {
	nt := NullableTime{}
	if string(data) == "null" || string(data) == "undefined" {
		nt.nilValue = true
	} else {
		var tt time.Time
		err := json.Unmarshal(data, &tt)
		if err != nil {
			return err
		}
		nt.t = tt
	}
	*t = nt
	return nil
}

func (t *NullableTime) Time() *time.Time {
	if t.nilValue {
		return nil
	}
	return &t.t
}

func (t *NullableTime) SetTime(time time.Time) {
	t.nilValue = false
	t.t = time
}

func (t *NullableTime) IsZero() bool {
	return t.t.IsZero()
}

func (t *NullableTime) IsNil() bool {
	return t.nilValue
}

func (t *NullableTime) SetNil() {
	t.nilValue = true
}

func (t *NullableTime) String() string {
	if t.nilValue {
		return "nil"
	}
	return t.t.String()
}

func (t *NullableTime) Format(layout string) string {
	if t.nilValue {
		return "nil"
	}
	return t.t.Format(layout)
}
