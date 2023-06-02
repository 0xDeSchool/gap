package ddd

import (
	"context"
	"errors"

	"github.com/0xDeSchool/gap/x"
)

type RepositoryBase[TEntity any, TKey comparable] interface {
	GetPagedList(ctx context.Context, p *x.PageAndSort) (*x.PagedResult[TEntity], error)
	GetById(ctx context.Context, id string) (*TEntity, error)
	GetOrNilById(ctx context.Context, id TKey) (*TEntity, error)
	GetMany(ctx context.Context, ids []TKey) ([]TEntity, error)
	Exists(ctx context.Context, id TKey) (bool, error)
	Count(ctx context.Context) (int64, error)
	FindByRegex(ctx context.Context, field, regex string, p *x.PageAndSort) (*x.PagedResult[TEntity], error)

	Insert(ctx context.Context, entity *TEntity) (*TEntity, error)
	// InsertMany ignoreErr 是否忽略批量插入时的错误, 一般为false, 当导入时忽略重复key的时候可以设为true
	InsertMany(ctx context.Context, entities []TEntity, ignoreErr bool) ([]TEntity, error)

	UpdateById(ctx context.Context, id TKey, data *TEntity) (int, error)

	Delete(ctx context.Context, id TKey) (int, error)
	DeleteMany(ctx context.Context, ids []TKey) (int, error)
}

var (
	ErrEntityNotFound = errors.New("entity not found")
)
