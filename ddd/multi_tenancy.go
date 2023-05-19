package ddd

type IMultiTenancy interface {
	GetTenant() string
	SetTenant(tenantId string)
}

type MultiTenantEntity struct {
	TenantId string `bson:"tenant_id"`
}

func (m *MultiTenantEntity) GetTenant() string {
	return m.TenantId
}

func (m *MultiTenantEntity) SetTenant(tenantId string) {
	m.TenantId = tenantId
}
