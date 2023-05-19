package multi_tenancy

import "github.com/0xDeSchool/gap/app"

func UseMultiTenancy(ab *app.AppBuilder) {
	ab.ConfigureServices(func() error {
		app.TryAddSingletonDefault[TenantResolver]()
		return nil
	})
}
