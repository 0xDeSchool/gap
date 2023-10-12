package aws

import (
	"github.com/0xDeSchool/gap/app"

	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func init() {
	app.Configure(func() error {
		opts := &Options{}
		utils.ViperBind("AWS", opts)
		app.AddValue(opts)
		app.AddSingleton(func() *aws.Config {
			conf, err := loadAwsConfig(opts)
			errx.CheckError(err)
			return conf
		})
		addS3Client()
		return nil
	})
}
