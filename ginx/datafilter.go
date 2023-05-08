package ginx

import (
	"context"
	"reflect"

	"github.com/gin-gonic/gin"
)

type DataFilterKey = string

const datafilterKey DataFilterKey = "_datafilter"

type GinDataFilterProvider struct {
}

type dataFiltersMap map[reflect.Type]struct{}

func (GinDataFilterProvider) DisableFilter(ctx context.Context, t reflect.Type) context.Context {
	filters := ctx.Value(datafilterKey)
	if filters == nil {
		filters = dataFiltersMap{}
		if gc, ok := ctx.(*gin.Context); ok {
			gc.Set(datafilterKey, filters)
		} else {
			ctx = context.WithValue(ctx, datafilterKey, filters)
		}
	}
	if f, ok := filters.(dataFiltersMap); ok {
		f[t] = struct{}{}
	} else {
		panic("invalid data filter")
	}
	return ctx
}

func (GinDataFilterProvider) EnableFilter(ctx context.Context, t reflect.Type) {
	filters := ctx.Value(datafilterKey)
	if filters == nil {
		return
	}
	if f, ok := filters.(dataFiltersMap); ok {
		delete(f, t)
	}
}

func (GinDataFilterProvider) IsFilterDisabled(ctx context.Context, t reflect.Type) bool {
	filters := ctx.Value(datafilterKey)
	if filters == nil {
		return false
	}
	if f, ok := filters.(dataFiltersMap); ok {
		_, ok := f[t]
		return ok
	}
	return false
}
