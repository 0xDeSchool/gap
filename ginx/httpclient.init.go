package ginx

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/utils"
)

func init() {
	app.Configure(func() error {
		var httpOptions = &HttpClientOptions{}
		utils.ViperBind("HttpClient", httpOptions)
		app.TryAddValue(httpOptions)
		app.TryAddTransient(func() *RequestClient {
			options := app.Get[HttpClientOptions]()
			return NewRequestClient(options)
		})
		return nil
	})
}
