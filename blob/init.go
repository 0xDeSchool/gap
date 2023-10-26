package blob

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/utils"
)

func init() {
	app.Configure(func() error {
		opts := NewOptions()
		app.AddValue(opts)
		utils.ViperBind("BlobStorage", opts)
		app.AddSingleton(NewBlobContainer)
		return nil
	})
}
