package ddd

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/x"
	"time"

	"github.com/0xDeSchool/gap/ginx"
)

type CreationAuditedEntity interface {
	Creating(ctx context.Context)
}

type UpdationAuditedEntity interface {
	Updating(ctx context.Context)
}

type ISoftDeleteEntity interface {
	Entity
	Deleting(ctx context.Context)
}

type AuditEntityBase struct {
	ID        string    `bson:"_id,omitempty"`
	CreatorId string    `bson:"creatorId"`
	CreatedAt time.Time `bson:"createdAt"`
}

func (e AuditEntityBase) GetId() string {
	return e.ID
}

func (e *AuditEntityBase) Creating(ctx context.Context) {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	if e.CreatorId == "" {
		e.CreatorId = ginx.CurrentUser(ctx).ID
	}
}

type FullAuditEntityBase struct {
	AuditEntityBase `bson:",inline"`
	UpdatedAt       time.Time `bson:"updatedAt,omitempty"`
	UpdaterId       string    `bson:"updaterId,omitempty"`
}

func (e *FullAuditEntityBase) Updating(ctx context.Context) {
	if e.UpdatedAt.IsZero() {
		e.UpdatedAt = time.Now()
	}
	if e.UpdaterId == "" {
		e.UpdaterId = ginx.CurrentUser(ctx).ID
	}
}

const SoftDeleteFieldName = "isDeleted"

type SoftDeleteEntity struct {
	IsDeleted  bool      `bson:"isDeleted"`
	DeletionAt time.Time `bson:"deletedAt,omitempty"`
	DeleterId  string    `bson:"deleterId,omitempty"`
}

func (e *SoftDeleteEntity) Deleting(ctx context.Context) {
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

func SetAudited(ctx context.Context, e any) {
	ig := *app.Get[x.IdGenerator[string]]()
	if ae, ok := e.(Entity); ok {
		if ae.GetId() == "" {
			ae.SetId(ig.Create())
		}
	}
	if ae, ok := e.(CreationAuditedEntity); ok {
		ae.Creating(ctx)
	}
	if ue, ok := e.(UpdationAuditedEntity); ok {
		ue.Updating(ctx)
	}

	if te, ok := e.(IMultiTenancy); ok {
		te.SetTenant(ginx.CurrentTenant(ctx).Id)
	}
}

// SetAuditedMany is a helper function to set audit fields for a slice of pointers.
// T is not struct ptr
func SetAuditedMany[T any](ctx context.Context, data []T) []any {
	result := make([]any, len(data))
	for i := range data {
		v := &data[i]
		SetAudited(ctx, v)
		result[i] = v
	}
	return result
}

// SetAuditedManyPtr is a helper function to set audit fields for a slice of pointers.
// T is struct ptr
func SetAuditedManyPtr[T any](ctx context.Context, data []T) []any {
	result := make([]any, len(data))
	for i := range data {
		v := data[i]
		SetAudited(ctx, v)
		result[i] = v
	}
	return result
}
