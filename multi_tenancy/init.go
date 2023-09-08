package multi_tenancy

import "github.com/0xDeSchool/gap/app"

func init() {
	app.Configure(func() error {
		app.TryAddSingletonDefault[TenantResolver]()
		return nil
	})
}
