package data

import (
	"context"
	"github.com/0xDeSchool/gap/x"
	"github.com/graph-gophers/dataloader/v7"
)

type ILoaderFactory[TKey comparable, T any] interface {
	Name() string
	CreateLoader() *dataloader.Loader[TKey, T]
}

type LoaderFactory[TKey comparable, T any] struct {
	batchFunc TypedBatchFunc[TKey, T]
	name      string
}

type TypedBatchFunc[TKey, T any] func(context.Context, []TKey) []T

func (tl *LoaderFactory[TKey, T]) Name() string {
	return tl.name
}

func (tl *LoaderFactory[TKey, T]) CreateLoader() *dataloader.Loader[TKey, T] {
	opt := dataloader.WithCache[TKey, T](NewContextCache[TKey, T](tl.name))
	return dataloader.NewBatchedLoader(tl.Load, opt)
}

func (tl *LoaderFactory[TKey, T]) Load(ctx context.Context, keys []TKey) []*dataloader.Result[T] {
	data := tl.batchFunc(ctx, keys)
	if len(data) != len(keys) {
		panic("len(output) != len(keys)")
	}
	output := make([]*dataloader.Result[T], len(keys))
	for index, v := range data {
		output[index] = &dataloader.Result[T]{Data: v, Error: nil}
	}
	return output
}

func NewLoaderWithName[TKey comparable, TResult any](name string, batchFunc TypedBatchFunc[TKey, TResult]) *LoaderFactory[TKey, TResult] {
	if name == "" {
		panic("name is empty")
	}
	tl := &LoaderFactory[TKey, TResult]{
		batchFunc: batchFunc,
		name:      name,
	}
	var _ ILoaderFactory[TKey, TResult] = tl
	return tl
}

func NewTypedLoader[TKey comparable, TResult any](batchFunc TypedBatchFunc[TKey, TResult]) *LoaderFactory[TKey, TResult] {
	return NewLoaderWithName[TKey, TResult](x.TypeOf[TResult]().Name(), batchFunc)
}
