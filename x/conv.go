package x

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func ToJsonMapString(v any) map[string]string {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string]string); ok {
		return m
	}

	fv := reflect.ValueOf(v)
	for fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
	}
	if fv.Kind() != reflect.Struct {
		panic("v must be struct")
	}
	c, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	m := map[string]any{}
	err = json.Unmarshal(c, &m)
	if err != nil {
		panic(err)
	}

	r := map[string]string{}
	for k, v := range m {
		r[k] = fmt.Sprintf("%v", v)
	}
	return r
}
