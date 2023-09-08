package store

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/ginx"
)

func init() {
	app.Configure(func() error {
		var p DataFilterProvider = &ginx.GinDataFilterProvider{}
		app.AddSingleton(func() *DataFilterProvider { return &p })
		return nil
	})
}
