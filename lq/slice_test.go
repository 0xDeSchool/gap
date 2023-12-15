package lq

import (
	"github.com/0xDeSchool/gap/x"
	"testing"
)

func TestSlice_Reduce(t *testing.T) {
	d := []int{1, 2, 3}

	x.AssertPanic(t, func() {
		First(d, nil)
	}, "parameter predicate is nil")

	res := Reduce(d, 0, func(a, b int) int {
		return a + b
	})
	sum := 6
	if res != 6 {
		t.Errorf("Reduce() = %v, want %v", res, sum)
	}

	res = Reduce(d, 0, func(a, b int) int {
		return a + b*2
	})
	sum = 12
	if res != 12 {
		t.Errorf("Reduce() = %v, want %v", res, sum)
	}
}

func TestSlice_First(t *testing.T) {
	res, has := First(nil, func(a int) bool { return true })
	if has || res != 0 {
		t.Errorf("First() = %v, want %v", has, false)
	}

	x.AssertPanic(t, func() {
		First([]int{1, 2, 3}, nil)
	}, "parameter predicate is nil")

	d := []int{1, 2, 3}
	res, has = First(d, func(a int) bool {
		return a == 2
	})
	if !has {
		t.Errorf("First() = %v, want %v", has, true)
	}
	if res != 2 {
		t.Errorf("First() = %v, want %v", res, 2)
	}

	res, has = First(d, func(a int) bool {
		return a == 6
	})
	if has {
		t.Errorf("First() = %v, want %v", has, false)
	}
	if res != 0 {
		t.Errorf("First() = %v, want %v", res, 0)
	}
}

func TestSlice_ToSet(t *testing.T) {
	d := []int{1, 2, 3}

	x.AssertPanic(t, func() {
		First(d, nil)
	}, "parameter predicate is nil")

	res := ToSet(d, func(i int) int {
		return i
	})
	if len(res) != 3 {
		t.Errorf("ToSet() = %v, want %v", len(res), 3)
	}
	if !res.Contains(1) {
		t.Errorf("ToSet() = %v, want %v", res, 1)
	}
	if !res.Contains(2) {
		t.Errorf("ToSet() = %v, want %v", res, 2)
	}
	if !res.Contains(3) {
		t.Errorf("ToSet() = %v, want %v", res, 3)
	}
	if res.Contains(4) {
		t.Errorf("ToSet() = %v, want %v", res, 4)
	}
}
