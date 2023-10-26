package s3

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/blob"
	"github.com/0xDeSchool/gap/x"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"net/url"
	"path"
	"strings"
)

func NewS3BlobProvider() blob.Provider {
	return &s3BlobProvider{}
}

type s3BlobProvider struct {
}

func (sp *s3BlobProvider) Save(ctx context.Context, b *blob.CreateBlobInfo) (*blob.SaveResult, error) {
	c := app.Get[s3.Client]()
	var bucket *string = nil
	opts := app.Get[s3BlobOptions]()
	if opts.Bucket() != "" {
		bucket = x.Ptr(opts.Bucket())
	}
	key := b.Name
	if !b.KeepName {
		key = sp.calculateKey(b)
	}
	input := &s3.PutObjectInput{
		Bucket:        bucket,
		Key:           &key,
		Body:          b.Blob,
		ContentType:   &b.ContentType,
		ContentLength: b.Size,
		//ContentDisposition: x.Ptr("attachment; filename=" + b.Name),
		Metadata: b.Metadata,
	}
	_, err := c.PutObject(ctx, input)
	if err != nil {
		return nil, err
	}
	u, err := url.JoinPath(opts.BaseUrl(), key)
	if err != nil {
		return nil, err
	}
	res := &blob.SaveResult{
		Url: u,
	}
	return res, nil
}

func (sp *s3BlobProvider) calculateKey(b *blob.CreateBlobInfo) string {
	ig := app.GetPtr[x.IdGenerator[string]]()
	ext := path.Ext(b.Name)
	p := "blob/" + x.Letters(2) + "/" + ig.Create()
	p = strings.ToLower(p) + ext
	return p
}
