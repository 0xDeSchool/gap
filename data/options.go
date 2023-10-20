package data

import (
	"github.com/0xDeSchool/gap/x"
	"reflect"
)

type LoaderOptions struct {
	Loaders map[string]any
}

func (lo *LoaderOptions) GetLoader(name string) any {
	if v, ok := lo.Loaders[name]; ok {
		return v
	}
	return nil
}

func NewLoaderOptions() *LoaderOptions {
	return &LoaderOptions{
		Loaders: map[string]any{},
	}
}

func AddLoader[TKey comparable, T any](opts *LoaderOptions, loader ILoaderFactory[TKey, *T]) {
	if loader == nil {
		panic("loader is nil")
	}
	name := loader.Name()
	if name == "" {
		name = reflect.TypeOf(loader).String()
	}
	opts.Loaders[name] = loader
}

func AddLoaderFunc[TKey comparable, T any](opts *LoaderOptions, batchFunc TypedBatchFunc[TKey, T]) {
	name := x.TypeOf[LoaderFactory[TKey, T]]().String()
	opts.Loaders[name] = NewLoaderWithName[TKey, T](name, batchFunc)
}

func AddNamedLoaderFunc[TKey comparable, T any](opts *LoaderOptions, name string, batchFunc TypedBatchFunc[TKey, T]) {
	if name == "" {
		panic("name is empty")
	}
	opts.Loaders[name] = NewLoaderWithName[TKey, T](name, batchFunc)
}
