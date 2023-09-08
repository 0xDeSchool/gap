package x

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/utils"
)

func init() {
	app.Configure(func() error {
		opts := &IdGeneratorOptions{
			WorkerId: 1,
		}
		utils.ViperBind("IdGenerator", opts)
		app.TryAddValue(opts)
		app.TryAddSingleton(NewStringIdGenerator[string])
		app.TryAddSingleton(NewNumberIdGenerator[int64])
		return nil
	})
}
