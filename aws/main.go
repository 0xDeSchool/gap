package aws

import (
	"context"
	"github.com/0xDeSchool/gap/app"

	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func AWS_SDK(b *app.AppBuilder) {
	b.ConfigureServices(func() error {
		opts := &AwsOptions{
			Url: "https://deschool.s3.amazonaws.com",
		}
		utils.ViperBind("AWS", opts)
		app.AddValue(opts)
		app.AddSingleton(func() *aws.Config {
			conf, err := LoadAwsConfig(opts)
			errx.CheckError(err)
			return conf
		})
		app.AddSingleton(func() *VideoTranscoder {
			config := app.Get[aws.Config]()
			return NewVideoTranscoder(config, opts)
		})
		return nil
	})
}

func LoadAwsConfig(opts *AwsOptions) (*aws.Config, error) {
	p := credentials.NewStaticCredentialsProvider(opts.AccessKeyId, opts.SecretAccessKey, opts.SessionToken)
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(opts.Region),
		config.WithCredentialsProvider(p),
	)
	return &cfg, err
}
