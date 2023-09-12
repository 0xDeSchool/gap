package multi_tenancy

import (
	"context"
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

type TenantResolver struct {
	// Name of resolver
	Name string
	// ResolveFunc tenant
	ResolveFunc func(ctx *TenantResolveContext) error
}

type TenantInfo struct {
	Id   string
	Name string
}

func (t *TenantInfo) IsHost() bool {
	return t.Id == ""
}
