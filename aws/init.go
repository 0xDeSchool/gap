package aws

import (
	"context"
	"github.com/0xDeSchool/gap/app"

	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func init() {
	app.Configure(func() error {
		opts := &AwsOptions{
			Url: "https://deschool.s3.amazonaws.com",
		}
		utils.ViperBind("AWS", opts)
		app.AddValue(opts)
		app.AddSingleton(func() *aws.Config {
			conf, err := loadAwsConfig(opts)
			errx.CheckError(err)
			return conf
		})
		app.AddSingleton(func() *VideoTranscoder {
			conf := app.Get[aws.Config]()
			return NewVideoTranscoder(conf, opts)
		})
		addS3Client()
		return nil
	})
}

func loadAwsConfig(opts *AwsOptions) (*aws.Config, error) {
	p := credentials.NewStaticCredentialsProvider(opts.AccessKeyId, opts.SecretAccessKey, opts.SessionToken)
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(opts.Region),
		config.WithCredentialsProvider(p),
	)
	return &cfg, err
}

func addS3Client() {
	app.TryAddSingleton(func() *s3.Client {
		conf := app.Get[aws.Config]()
		return s3.NewFromConfig(*conf)
	})
}
