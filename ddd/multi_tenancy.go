package ddd

type IMultiTenancy interface {
	GetTenant() string
	SetTenant(tenantId string)
}

const TenantIdDbKey = "tenantId"

type MultiTenantEntity struct {
	TenantId string `bson:"tenantId"`
}

func (m *MultiTenantEntity) GetTenant() string {
	return m.TenantId
}

func (m *MultiTenantEntity) SetTenant(tenantId string) {
	m.TenantId = tenantId
}
