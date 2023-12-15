package x

import "reflect"

func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func HasField(t reflect.Type, name string) bool {
	_, has := t.FieldByName(name)
	return has
}

func Zero[T any]() T {
	var res T
	return res
}
