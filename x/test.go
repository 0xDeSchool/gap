package x

import "testing"

func AssertPanic(t *testing.T, f func(), msg interface{}) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic: %v", msg)
		}
	}()
	f()
}
