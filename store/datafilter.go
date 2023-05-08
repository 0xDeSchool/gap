package store

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"reflect"

	"github.com/0xDeSchool/gap/x"
)

type Datafilter interface {
	Filter(ctx context.Context, v any) any
}

type DataFilterProvider interface {
	DisableFilter(ctx context.Context, t reflect.Type) context.Context
	EnableFilter(ctx context.Context, t reflect.Type)
	IsFilterDisabled(ctx context.Context, t reflect.Type) bool
}

// type dataFiltersMap map[reflect.Type]struct{}

// type dataFilterKeyType struct {
// }

// DisableFilter disables data filter for TEntity
func DisableFilter[TEntity any](ctx context.Context) context.Context {
	p := *app.Get[DataFilterProvider]()
	return p.DisableFilter(ctx, x.TypeOf[TEntity]())
}

// EnableFilter enables data filter for TEntity
func EnableFilter[TEntity any](ctx context.Context) {
	p := *app.Get[DataFilterProvider]()
	p.EnableFilter(ctx, x.TypeOf[TEntity]())
}

// IsFilterDisabled returns true if data filter is disabled for TEntity
func IsFilterDisabled[TEntity any](ctx context.Context) bool {
	p := *app.Get[DataFilterProvider]()
	return p.IsFilterDisabled(ctx, x.TypeOf[TEntity]())
}
