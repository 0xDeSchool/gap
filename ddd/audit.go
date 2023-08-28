package ddd

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/x"
	"time"

	"github.com/0xDeSchool/gap/ginx"
)

type CreationAuditedEntity[TKey comparable] interface {
	Creating(ctx context.Context)
}

type UpdationAuditedEntity[TKey comparable] interface {
	Updating(ctx context.Context)
}

type ISoftDeleteEntity[TKey comparable] interface {
	Entity[TKey]
	Deleting(ctx context.Context)
}

type AuditEntityBase[TKey comparable] struct {
	EntityBase[TKey] `bson:",inline"`
	CreatorId        string    `bson:"creatorId"`
	CreatedAt        time.Time `bson:"createdAt"`
}

func (e *AuditEntityBase[TKey]) GetId() TKey {
	return e.ID
}

func (e *AuditEntityBase[TKey]) Creating(ctx context.Context) {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	if e.CreatorId == "" {
		e.CreatorId = ginx.CurrentUser(ctx).ID
	}
}

type FullAuditEntityBase[TKey comparable] struct {
	AuditEntityBase[TKey] `bson:",inline"`
	UpdatedAt             time.Time `bson:"updatedAt,omitempty"`
	UpdaterId             string    `bson:"updaterId,omitempty"`
}

func (e *FullAuditEntityBase[TKey]) Updating(ctx context.Context) {
	if e.UpdatedAt.IsZero() {
		e.UpdatedAt = time.Now()
	}
	if e.UpdaterId == "" {
		e.UpdaterId = ginx.CurrentUser(ctx).ID
	}
}

const SoftDeleteFieldName = "isDeleted"

type SoftDeleteEntity[TKey comparable] struct {
	IsDeleted  bool      `bson:"isDeleted"`
	DeletionAt time.Time `bson:"deletedAt,omitempty"`
	DeleterId  string    `bson:"deleterId,omitempty"`
}

func (e *SoftDeleteEntity[TKey]) Deleting(ctx context.Context) {
	if e.DeletionAt.IsZero() {
		e.DeletionAt = time.Now()
	}
	if e.DeleterId == "" {
		e.DeleterId = ginx.CurrentUser(ctx).ID
	}
	e.IsDeleted = true
}

type hardDeleteKey struct{}

var hardKey = hardDeleteKey{}

func WithHardDelete(ctx context.Context) context.Context {
	return context.WithValue(ctx, hardKey, true)
}

func SetAudited[TKey comparable](ctx context.Context, e any) {
	if ae, ok := e.(Entity[TKey]); ok {
		if ig, ok2 := app.GetPtrOptional[x.IdGenerator[TKey]](); ok2 {
			var defaultId TKey
			if ae.GetId() == defaultId {
				ae.SetId(ig.Create())
			}
		}
	}
	if ae, ok := e.(CreationAuditedEntity[TKey]); ok {
		ae.Creating(ctx)
	}
	if ue, ok := e.(UpdationAuditedEntity[TKey]); ok {
		ue.Updating(ctx)
	}

	if te, ok := e.(IMultiTenancy); ok {
		te.SetTenant(ginx.CurrentTenant(ctx).Id)
	}
}

// SetAuditedMany is a helper function to set audit fields for a slice of pointers.
// T is not struct ptr
func SetAuditedMany[T any, TKey comparable](ctx context.Context, data []T) []any {
	result := make([]any, len(data))
	for i := range data {
		v := &data[i]
		SetAudited[TKey](ctx, v)
		result[i] = v
	}
	return result
}

// SetAuditedManyPtr is a helper function to set audit fields for a slice of pointers.
// T is struct ptr
func SetAuditedManyPtr[T any, TKey comparable](ctx context.Context, data []*T) []any {
	result := make([]any, len(data))
	for i := range data {
		v := data[i]
		SetAudited[TKey](ctx, v)
		result[i] = v
	}
	return result
}
