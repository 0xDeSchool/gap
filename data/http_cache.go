package data

import (
	"context"
	"fmt"
	"github.com/graph-gophers/dataloader/v7"
	"sync"
)

type cacheKeyType = string

const cacheKey cacheKeyType = "dataloader_cache"

type cacheDataStore struct {
	mu   sync.RWMutex
	data map[any]any
}

func (c *cacheDataStore) Set(key any, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *cacheDataStore) Get(key any) (value any, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok = c.data[key]
	return
}

func (c *cacheDataStore) Delete(key any) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ok {
		delete(c.data, key)
	}
	return
}

type InContextCache[TKey comparable, TValue any] struct {
	prefix string
}

// NewContextCache constructs a new InMemoryCache
func NewContextCache[TKey comparable, TValue any](name string) *InContextCache[TKey, TValue] {
	return &InContextCache[TKey, TValue]{
		prefix: name,
	}
}

// Set sets the `value` at `key` in the cache
func (c *InContextCache[TKey, TValue]) Set(ctx context.Context, key TKey, value dataloader.Thunk[TValue]) {
	cache := c.getCache(ctx)
	if cache != nil {
		cache.Set(c.key(key), value)
	}
}

// Get gets the value at `key` if it exists, returns value (or nil) and bool
// indicating of value was found
func (c *InContextCache[TKey, TValue]) Get(ctx context.Context, key TKey) (dataloader.Thunk[TValue], bool) {
	cache := c.getCache(ctx)
	if cache != nil {
		if v, ok := cache.Get(c.key(key)); ok {
			return v.(dataloader.Thunk[TValue]), true
		}
	}
	return nil, false
}

// Delete deletes item at `key` from cache
func (c *InContextCache[TKey, TValue]) Delete(ctx context.Context, key TKey) bool {
	cache := c.getCache(ctx)
	if cache != nil {
		return cache.Delete(c.key(key))
	}
	return false
}

// Clear clears the cache
func (c *InContextCache[TKey, TValue]) Clear() {
}

func (c *InContextCache[TKey, TValue]) getCache(ctx context.Context) (cache *cacheDataStore) {
	v := ctx.Value(cacheKey)
	if v == nil {
		cache = nil
	} else {
		cache = v.(*cacheDataStore)
	}
	return
}

func (c *InContextCache[TKey, TValue]) key(k TKey) string {
	return fmt.Sprintf("%s:%v", c.prefix, k)
}

func WithDataLoaderCache(ctx context.Context) context.Context {
	cache := &cacheDataStore{
		data: map[any]any{},
	}
	return context.WithValue(ctx, cacheKey, cache)
}
