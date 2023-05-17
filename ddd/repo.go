package ddd

import (
	"context"
	"errors"

	"github.com/0xDeSchool/gap/x"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositoryBase[TEntity any] interface {
	GetPagedList(ctx context.Context, p *x.PageAndSort) (*x.PagedResult[TEntity], error)
	GetById(ctx context.Context, id primitive.ObjectID) (*TEntity, error)
	GetOrNilById(ctx context.Context, id primitive.ObjectID) (*TEntity, error)
	GetMany(ctx context.Context, ids []primitive.ObjectID) ([]TEntity, error)
	Exists(ctx context.Context, id primitive.ObjectID) (bool, error)
	Count(ctx context.Context) (int64, error)
	FindByRegex(ctx context.Context, field, regex string, p *x.PageAndSort) (*x.PagedResult[TEntity], error)

	Insert(ctx context.Context, entity *TEntity) (primitive.ObjectID, error)
	// InsertMany ignoreErr 是否忽略批量插入时的错误, 一般为false, 当导入时忽略重复key的时候可以设为true
	InsertMany(ctx context.Context, entities []TEntity, ignoreErr bool) ([]primitive.ObjectID, error)

	UpdateById(ctx context.Context, id primitive.ObjectID, data *TEntity) (int, error)

	Delete(ctx context.Context, id primitive.ObjectID) (int, error)
	DeleteMany(ctx context.Context, ids []primitive.ObjectID) (int, error)
}

var (
	ErrEntityNotFound = errors.New("entity not found")
)
