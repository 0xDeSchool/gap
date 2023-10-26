package blob

import "context"

type Provider interface {
	Save(ctx context.Context, b *CreateBlobInfo) (*SaveResult, error)
}
