package x

import (
	"strconv"

	"github.com/0xDeSchool/gap/errx"
)

func ToUint(s string) uint64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseUint(s, 10, 64)
	errx.CheckError(err)
	return v
}

// CanConvert returns true if FromType can be converted to ToType.
func CanConvert[FromType any, ToType any]() bool {
	var v any = (*FromType)(nil)
	if _, ok := v.(ToType); ok {
		return ok
	}
	return false
}

func Ptr[T any](v T) *T {
	return &v
}

func Ptrs[T any](v []T) []*T {
	r := make([]*T, 0, len(v))
	for _, v := range v {
		r = append(r, &v)
	}
	return r
}
