package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type Options struct {
	BaseUrl         string
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}

func loadAwsConfig(opts *Options) (*aws.Config, error) {
	p := credentials.NewStaticCredentialsProvider(opts.AccessKeyId, opts.SecretAccessKey, opts.SessionToken)
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(opts.Region),
		config.WithCredentialsProvider(p),
	)
	return &cfg, err
}
