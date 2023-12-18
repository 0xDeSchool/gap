package x

import (
	"testing"
	"time"
)

// TestToMapString  test ToJsonMapString method
func TestToMapString(t *testing.T) {
	v1 := "test"
	AssertPanic(t, func() {
		ToJsonMapString(v1)
	}, "v must be struct")

	v2 := struct {
		Name     string `json:"name"`
		Age      int    `json:"age"`
		Birthday *time.Time
	}{Name: "test", Age: 12, Birthday: Ptr(time.Now())}
	m2 := ToJsonMapString(v2)
	if m2["name"] != "test" {
		t.Errorf("ToJsonMapString() Name = %v, want %v", m2["Name"], "test")
	}

	if m2["age"] != "12" {
		t.Errorf("ToJsonMapString() Age = %v, want %v", m2["Age"], "12")
	}
	if m2["Birthday"] != v2.Birthday.Format(time.RFC3339Nano) {
		t.Errorf("ToJsonMapString() Birthday = %v, want %v", m2["Birthday"], v2.Birthday.String())
	}

}
