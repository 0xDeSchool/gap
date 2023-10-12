package aws

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Bucket s3 client options
func (opts Options) Bucket() string {
	return opts.GetString("Bucket")
}

// SetBucket s3 client options
func (opts Options) SetBucket(v string) {
	opts["Bucket"] = v
}

func addS3Client() {
	app.TryAddSingleton(func() *s3.Client {
		conf := app.Get[aws.Config]()
		return s3.NewFromConfig(*conf)
	})
}
