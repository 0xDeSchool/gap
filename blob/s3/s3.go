package s3

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/aws"
	"github.com/0xDeSchool/gap/blob"
)

const ProviderConfigKey = "s3"

type s3BlobOptions aws.Options

func (opts s3BlobOptions) Bucket() string {
	return aws.Options(opts).Bucket()
}

func (opts s3BlobOptions) BaseUrl() string {
	return aws.Options(opts).BaseUrl()
}

func newS3BlobOptions() *s3BlobOptions {
	awsOpts := app.Get[aws.Options]()
	blobOpts := app.Get[blob.Options]()
	res := make(s3BlobOptions)
	opts := aws.Options(res)
	for k, v := range awsOpts.Values() {
		opts.Set(k, v)
	}
	if cfg, ok := blobOpts.ProvidersConfig[ProviderConfigKey]; ok {
		for k, v := range cfg {
			opts.Set(k, v)
		}
	}
	return &res
}

func UseBlobS3Storage(providerKey blob.ProviderType) {
	app.AddSingleton(newS3BlobOptions)
	app.ConfigureOptions(func(c *app.Container, opts *blob.Options) {
		p := NewS3BlobProvider()
		opts.AddProvider(providerKey, p)
		if opts.DefaultProvider() == nil {
			opts.SetDefaultProvider(p)
		}
	})
}
