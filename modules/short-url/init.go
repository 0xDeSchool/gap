package short_url

import (
	"github.com/0xDeSchool/gap/app"
)

func init() {
	app.Configure(func() error {
		app.AddTransient(NewMongoRepository)
		app.AddSingleton(NewDefaultKeyCreator)
		return nil
	})
}
