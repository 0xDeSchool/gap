package multi_tenancy

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/cache"
)

type TenantResolveContext struct {
	Ctx            context.Context
	TenantIdOrName string
	Handled        bool
}

func NewTenantResolveContext(ctx context.Context) *TenantResolveContext {
	return &TenantResolveContext{}
}

func (c *TenantResolveContext) HasResolved() bool {
	return c.Handled || c.TenantIdOrName != ""
}

type ITenantResolveContributor interface {
	// Name of resolver
	Name() string
	// Resolve tenant
	Resolve(ctx *TenantResolveContext) error
}

type TenantInfo struct {
	Id   string
	Name string
}

func (t *TenantInfo) IsHost() bool {
	return t.Id == ""
}

type TenantResolver struct {
	cache cache.Cache
}

func NewTenantResolver() *TenantResolver {
	return &TenantResolver{
		cache: cache.New(&cache.CacheOptions{
			LifeWindow: 0,
		}),
	}
}

func (t *TenantResolver) ResolveTenantIdOrName(ctx context.Context) (*TenantInfo, error) {
	result := &TenantInfo{}
	resolvers := app.GetArray[ITenantResolveContributor]()
	resolveCtx := NewTenantResolveContext(ctx)
	for _, resolver := range resolvers {
		err := resolver.Resolve(resolveCtx)
		if err != nil {
			return nil, err
		}
		if resolveCtx.HasResolved() {
			// TODO: support tenant name
			result.Id = resolveCtx.TenantIdOrName
			break
		}
	}
	return result, nil
}
