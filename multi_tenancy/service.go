package multi_tenancy

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/cache"
)

type TenantService struct {
	cache cache.Cache

	opts *TenantOptions
}

func NewTenantService() *TenantService {
	return &TenantService{
		cache: cache.New(&cache.CacheOptions{
			LifeWindow: 0,
		}),
		opts: app.Get[TenantOptions](),
	}
}

func (t *TenantService) ResolveTenant(ctx context.Context) (*TenantInfo, error) {
	result := &TenantInfo{}
	resolveCtx := NewTenantResolveContext(ctx)
	for _, resolver := range t.opts.Resolvers {
		err := resolver.ResolveFunc(resolveCtx)
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
