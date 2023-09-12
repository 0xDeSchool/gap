package ginx

import "github.com/0xDeSchool/gap/app"

func init() {
	app.Configure(func() error {
		app.AddValue(&AuthOptions[string]{
			handlers: []AuthHandler[string]{},
		})
		return nil
	})
}
