package mongodb

import (
	"context"

	"github.com/0xDeSchool/gap/store"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoDataFilter func(ctx context.Context, filter bson.D) bson.D

type mongoDataFilter[T any] struct {
	filter MongoDataFilter
}

func (f *mongoDataFilter[T]) Filter(ctx context.Context, v any) any {
	return f.filter(ctx, v.(bson.D))
}

func newMongoDataFilter[T any](f MongoDataFilter) *mongoDataFilter[T] {
	return &mongoDataFilter[T]{
		filter: f,
	}
}

// AddFilter add filter for TEntity
func AddFilter[T any](opts *store.StoreOptions, filter MongoDataFilter) {
	store.AddFilter[T](opts, newMongoDataFilter[T](filter))
}
