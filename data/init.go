package data

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/ginx"
)

func init() {
	app.Configure(func() error {
		opts := NewLoaderOptions()
		app.AddValue(opts)
		return nil
	})
	ginx.PreConfigure(func(s *ginx.Server) error {
		s.Use(nil, dataLoaderFunc())
		return nil
	})
}
