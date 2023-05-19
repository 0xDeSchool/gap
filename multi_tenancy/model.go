package multi_tenancy

import "github.com/0xDeSchool/gap/ddd"

type Tenant struct {
	ddd.AuditEntityBase `bson:",inline"`
	Name                string `bson:"name"`
	DisplayName         string `bson:"display_name"`
}
