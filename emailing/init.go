package emailing

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/utils"
)

func init() {
	app.Configure(func() error {
		emailOpts := &EmailOptions{}
		utils.ViperBind("Email", emailOpts)
		app.AddValue(emailOpts)
		app.AddSingleton(func() *EmailSender {
			var es EmailSender = NewGomailSender(emailOpts)
			return &es
		})
		return nil
	})
}
