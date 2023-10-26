package blob

import (
	"context"
	"errors"
	"github.com/0xDeSchool/gap/app"
	"io"
)

var ErrProviderNotFound = errors.New("blob provider not found")

type CreateBlobInfo struct {
	Name         string
	Size         int64
	Blob         io.Reader
	ContentType  string
	ProviderType *ProviderType
	KeepName     bool // keep original name as key

	Metadata map[string]string
}

type Container struct {
	opts *Options
}

func NewBlobContainer() *Container {
	return &Container{
		opts: app.Get[Options](),
	}
}

func (c *Container) Save(ctx context.Context, b *CreateBlobInfo) (*SaveResult, error) {
	pt := defaultProviderKey
	if b.ProviderType != nil {
		pt = *b.ProviderType
	}
	p := c.opts.GetProvider(pt)
	if p == nil {
		return nil, ErrProviderNotFound
	}
	res, err := p.Save(ctx, b)
	if err != nil {
		return nil, err
	}
	e := &CreatedEventData{
		Path:   b.Name,
		Size:   b.Size,
		Result: res,
	}

	for _, handler := range c.opts.Handlers() {
		if err2 := handler(ctx, e); err2 != nil {
			return nil, err
		}
	}
	return res, nil
}
