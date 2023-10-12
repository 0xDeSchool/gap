package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"strings"
)

type Options map[string]any

func (opts Options) Set(k string, v any) {
	opts[strings.ToLower(k)] = v
}

func (opts Options) Get(k string) any {
	return opts[strings.ToLower(k)]
}

func (opts Options) Values() map[string]any {
	return opts
}

func (opts Options) GetString(k string) string {
	if v, ok := opts[strings.ToLower(k)].(string); ok {
		return v
	}
	return ""
}

func (opts Options) BaseUrl() string {
	return opts.GetString("BaseUrl")
}
func (opts Options) SetBaseUrl(v string) {
	opts.Set("BaseUrl", v)
}

func (opts Options) AccessKeyId() string {
	return opts.GetString("AccessKeyId")
}
func (opts Options) SetAccessKeyId(v string) {
	opts.Set("AccessKeyId", v)
}

func (opts Options) SecretAccessKey() string {
	return opts.GetString("SecretAccessKey")
}
func (opts Options) SetSecretAccessKey(v string) {
	opts.Set("SecretAccessKey", v)
}

func (opts Options) SessionToken() string {
	return opts.GetString("SessionToken")
}
func (opts Options) SetSessionToken(v string) {
	opts.Set("SessionToken", v)
}

func (opts Options) Region() string {
	return opts.GetString("Region")
}
func (opts Options) SetRegion(v string) {
	opts.Set("Region", v)
}

func loadAwsConfig(opts *Options) (*aws.Config, error) {
	p := credentials.NewStaticCredentialsProvider(opts.AccessKeyId(), opts.SecretAccessKey(), opts.SessionToken())
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(opts.Region()),
		config.WithCredentialsProvider(p),
	)
	return &cfg, err
}
