package multi_tenancy

import "context"

var TenantKey = "Tenant"
var TenantFilterKey = "TenantFilter"

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
