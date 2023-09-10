package ginx

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/utils"
)

var builder = NewServerBuilder()

func init() {
	builder.Options.Port = 5000
	app.Configure(func() error {
		utils.ViperBind("Server", builder.Options)
		app.AddValue(builder.Options)
		return nil
	})
	app.OrderRun(999, func() error {
		s, err := builder.Build()
		if err != nil {
			return err
		}
		return s.Run()
	})
}

// PreConfigure 配置服务，在Configure之前运行
func PreConfigure(action ServerConfigureFunc) {
	builder.PreConfigure(action)
}

// Configure 配置服务，在 gin.Run 之前运行
func Configure(action ServerConfigureFunc) {
	builder.Configure(action)
}
