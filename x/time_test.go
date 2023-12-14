package x

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNullableTime_IsZero(t *testing.T) {
	nt := NullableTime{}
	if !nt.IsZero() {
		t.Errorf("NullableTime.IsZero() = %v, want %v", nt.IsZero(), true)
	}
}

func TestNullableTime_IsNil(t *testing.T) {
	nt := NullableTime{}
	if nt.IsNil() {
		t.Errorf("NullableTime.IsNil() = %v, want %v", nt.IsNil(), true)
	}
}

func TestNullableTime_SetNil(t *testing.T) {
	nt := NullableTime{}
	nt.SetNil()
	if !nt.IsNil() {
		t.Errorf("NullableTime.SetNil() = %v, want %v", nt.IsNil(), true)
	}
}

func TestNullableTime_SetTime(t *testing.T) {
	now := time.Now()
	nt := NullableTime{}
	nt.SetTime(now)
	if nt.IsNil() {
		t.Errorf("NullableTime.SetTime() = %v, want %v", nt.IsNil(), false)
	}
	if *nt.Time() != now {
		t.Errorf("NullableTime.SetTime() = %v, want %v", *nt.Time(), now)
	}
}

func TestNullableTime_String(t *testing.T) {
	now := time.Now()
	nt := NullableTime{}
	nt.SetTime(now)
	if nt.String() != now.String() {
		t.Errorf("NullableTime.String() = %v, want %v", nt.String(), now.String())
	}
	nt.SetNil()
	if nt.String() != "nil" {
		t.Errorf("NullableTime.String() = %v, want %v", nt.String(), "nil")
	}
}

func TestNullableTime_Format(t *testing.T) {
	now := time.Now()
	nt := NullableTime{}
	nt.SetTime(now)
	if nt.Format("2006-01-02") != now.Format("2006-01-02") {
		t.Errorf("NullableTime.Format() = %v, want %v", nt.Format("2006-01-02"), now.Format("2006-01-02"))
	}
	nt.SetNil()
	if nt.Format("2006-01-02") != "nil" {
		t.Errorf("NullableTime.Format() = %v, want %v", nt.Format("2006-01-02"), "nil")
	}
}

func TestNullableTime_UnmarshalJSON(t *testing.T) {
	var a struct {
		Time NullableTime
	}
	jc := `{"Time":"2019-01-01T00:00:00Z"}`
	_ = json.Unmarshal([]byte(jc), &a)
	if a.Time.Time().Year() != 2019 {
		t.Errorf("NullableTime.UnmarshalJSON() = %v, want %v", a.Time.Time().Year(), 2019)
	}

	jc = `{"Time":null}`
	_ = json.Unmarshal([]byte(jc), &a)
	if !a.Time.IsNil() {
		t.Errorf("NullableTime.UnmarshalJSON() = %v, want %v", a.Time.IsNil(), true)
	}

	jc = `{}`
	_ = json.Unmarshal([]byte(jc), &a)
	if !a.Time.IsZero() { // zero value of time.Time is 0001-01-01 00:00:00 +0000 UTC
		t.Errorf("NullableTime.UnmarshalJSON() = %v, want %v", a.Time.IsZero(), true)
	}
}
