package multi_tenancy

import (
	"context"
	"github.com/0xDeSchool/gap/app"
)

type TenantResolveContext struct {
	Ctx            context.Context
	TenantIdOrName string
	Handled        bool
}

func NewTenantResolveContext(ctx context.Context) *TenantResolveContext {
	return &TenantResolveContext{
		Ctx: ctx,
	}
}

func (c *TenantResolveContext) HasResolved() bool {
	return c.Handled || c.TenantIdOrName != ""
}

type TenantOptions struct {
	Resolvers []TenantResolver
}

func NewTenantOptions() *TenantOptions {
	return &TenantOptions{
		Resolvers: make([]TenantResolver, 0),
	}
}

type ResolveTenantFunc func(ctx *TenantResolveContext) error

type TenantResolver struct {
	// Name of resolver
	Name string
	// ResolveFunc tenant
	ResolveFunc ResolveTenantFunc
}

type TenantInfo struct {
	Id   string
	Name string
}

func (t *TenantInfo) IsHost() bool {
	return t.Id == ""
}

func AddResolver(name string, h ResolveTenantFunc) {
	app.ConfigureOptions(func(c *app.Container, opts *TenantOptions) {
		opts.Resolvers = append(opts.Resolvers, TenantResolver{
			Name:        name,
			ResolveFunc: h,
		})
	})
}
