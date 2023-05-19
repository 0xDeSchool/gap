package multi_tenancy

import "context"

var TenantKey = struct{}{}
var TenantFilterKey = struct{}{}

type ITenantStore interface {
	// FindById Get tenant by id
	FindById(tenantId string) (*Tenant, error)
	// FindByName Get tenant by name
	FindByName(name string) (*Tenant, error)
}

func WithTenant(ctx context.Context, tenant *TenantInfo) context.Context {
	return context.WithValue(ctx, TenantKey, tenant)
}

func CurrentTenant(ctx context.Context) *TenantInfo {
	v := ctx.Value(TenantKey)
	if v == nil {
		return &TenantInfo{}
	}
	return v.(*TenantInfo)
}

func DisableMultiTenantFilter(ctx context.Context) context.Context {
	return context.WithValue(ctx, TenantKey, false)
}

func IsEnableMultiTenantFilter(ctx context.Context) bool {
	if v, ok := ctx.Value(TenantKey).(bool); ok {
		return v
	}
	return true
}
