package ginx

import "github.com/0xDeSchool/gap/app"

func init() {
	app.Configure(func() error {
		app.AddValue(&ResultHandlerOptions{})
		return nil
	})
}
