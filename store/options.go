package store

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"reflect"

	"github.com/0xDeSchool/gap/x"
)

type dataFilters map[reflect.Type]map[Datafilter]struct{}

type StoreOptions struct {
	dataFilters dataFilters
}

var golbalType = reflect.TypeOf(struct{}{})

func NewStoreOptions() *StoreOptions {
	return &StoreOptions{
		dataFilters: dataFilters{},
	}
}

// AddFilter add filter for TEntity
func AddFilter[T any](opts *StoreOptions, filter ...Datafilter) *StoreOptions {
	t := x.TypeOf[T]()
	for _, f := range filter {
		if _, ok := opts.dataFilters[t]; !ok {
			opts.dataFilters[t] = map[Datafilter]struct{}{}
		}
		opts.dataFilters[t][f] = struct{}{}
	}
	return opts
}

func AddGlobalFilter(opts *StoreOptions, filter Datafilter) *StoreOptions {
	if _, ok := opts.dataFilters[golbalType]; !ok {
		opts.dataFilters[golbalType] = map[Datafilter]struct{}{}
	}
	opts.dataFilters[golbalType][filter] = struct{}{}
	return opts
}

// DataFilters get filter for TEntity
func DataFilters[T any](ctx context.Context, opts *StoreOptions) []Datafilter {
	p := *app.Get[DataFilterProvider]()
	filters := make([]Datafilter, 0)
	if v, ok := opts.dataFilters[golbalType]; ok {
		for k, _ := range v {
			filters = append(filters, k)
		}
	}
	t := x.TypeOf[T]()
	if p.IsFilterDisabled(ctx, t) {
		return filters
	}
	if v, ok := opts.dataFilters[t]; ok {
		for k, _ := range v {
			filters = append(filters, k)
		}
	}
	return filters
}
