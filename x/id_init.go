package x

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/utils"
)

func init() {
	app.Configure(func() error {
		opts := &IdGeneratorOptions{}
		utils.ViperBind("IdGenerator", opts)
		app.AddValue(opts)
		app.AddSingleton(NewStringIdGenerator[string])
		app.AddSingleton(NewNumberIdGenerator[int64])
		return nil
	})
}
