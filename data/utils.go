package data

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/x"
	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/dataloader/v7"
	"sync"
)

type ctxLoaders struct {
	loaders map[string]any
	mutex   sync.Mutex
}

func (c *ctxLoaders) Get(name string) (any, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, ok := c.loaders[name]
	return v, ok
}

func (c *ctxLoaders) Set(name string, loader any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.loaders[name] = loader
}

func (c *ctxLoaders) GetOrSet(name string, loader any) any {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.loaders[name]; ok {
		return v
	} else {
		c.loaders[name] = loader
		return loader
	}
}

type loaderKeyType struct{}

var loaderKey loaderKeyType = struct{}{}

func GetTypedLoader[TKey comparable, T any](ctx context.Context) *dataloader.Loader[TKey, T] {
	name := x.TypeOf[LoaderFactory[TKey, T]]().String()
	return GetLoader[TKey, T](ctx, name)
}

func GetLoader[TKey comparable, T any](ctx context.Context, name string) *dataloader.Loader[TKey, T] {
	if v, ok := ctx.Value(loaderKey).(*ctxLoaders); ok {
		if loader, ok2 := v.Get(name); ok2 {
			return loader.(*dataloader.Loader[TKey, T])
		} else {
			loader = createLoader[TKey, T](name)
			loader = v.GetOrSet(name, loader)
			return loader.(*dataloader.Loader[TKey, T])
		}
	}
	panic("loader not found")
}

func createLoader[TKey comparable, T any](name string) *dataloader.Loader[TKey, T] {
	lm := app.Get[LoaderOptions]()
	if v := lm.GetLoader(name); v != nil {
		if l, ok := v.(ILoaderFactory[TKey, T]); ok {
			return l.CreateLoader()
		}
		panic("loader type error")
	}
	panic("loader not found")
}

func dataLoaderFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		newCtx := context.WithValue(c.Request.Context(), loaderKey, &ctxLoaders{
			loaders: map[string]any{},
		})
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
