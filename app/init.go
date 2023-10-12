package app

import (
	"github.com/0xDeSchool/gap/utils"
	"github.com/rs/zerolog"
)

func init() {
	Configure(func() error {
		opts := &AppOptions{
			LogLevel: zerolog.InfoLevel,
		}
		utils.ViperBind("Application", opts)
		AddValue(opts)
		return nil
	})
}
