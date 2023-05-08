package mongodb

import (
	"context"
	"errors"

	"github.com/0xDeSchool/gap/ginx"
	"github.com/0xDeSchool/gap/log"
	"github.com/0xDeSchool/gap/x"

	"github.com/0xDeSchool/gap/ddd"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/store"
	"github.com/0xDeSchool/gap/utils/linq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[TEntity any] struct {
	c    *mongo.Collection
	opts *store.StoreOptions
}

func NewCollection[TEntity any](c *mongo.Collection, opts *store.StoreOptions) *Collection[TEntity] {
	return &Collection[TEntity]{
		c:    c,
		opts: opts,
	}
}

func (c *Collection[TEntity]) Find(ctx context.Context, filter bson.D, opts ...*options.FindOptions) []TEntity {
	filter = c.setSoftDeleteFilter(filter)
	filter = c.MergeGlobalFilter(ctx, filter)
	cur, err := c.Col().Find(ctx, filter, opts...)
	errx.CheckError(err)
	data := make([]TEntity, 0)
	errx.CheckError(cur.All(context.Background(), &data))
	return data
}

func (c *Collection[TEntity]) FindOne(ctx context.Context, filter bson.D, opts ...*options.FindOneOptions) *TEntity {
	filter = c.setSoftDeleteFilter(filter)
	filter = c.MergeGlobalFilter(ctx, filter)
	result := c.Col().FindOne(ctx, filter, opts...)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		errx.PanicNotFound("data not found")
	} else {
		errx.CheckError(result.Err())
	}
	var v TEntity
	err := result.Decode(&v)
	errx.CheckError(err)
	return &v
}

func (c *Collection[TEntity]) Get(ctx context.Context, id primitive.ObjectID) *TEntity {
	errx.NotNil(id, "id")
	filter := bson.D{{Key: "_id", Value: id}}
	filter = c.MergeGlobalFilter(ctx, filter)
	return c.FindOne(ctx, filter)
}

func (c *Collection[TEntity]) Count(ctx context.Context, filter bson.D, opts ...*options.CountOptions) int64 {
	filter = c.setSoftDeleteFilter(filter)
	filter = c.MergeGlobalFilter(ctx, filter)
	totalCount, err := c.Col().CountDocuments(ctx, filter, opts...)
	errx.CheckError(err)
	return totalCount
}

func (c *Collection[TEntity]) GetMany(ctx context.Context, ids []primitive.ObjectID, filter bson.D) []TEntity {
	errx.NotNil(ids, "ids")
	if len(ids) == 0 {
		return make([]TEntity, 0)
	}
	f := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	if len(filter) > 0 {
		f = append(f, filter...)
	}
	f = c.setSoftDeleteFilter(f)
	f = c.MergeGlobalFilter(ctx, f)
	result, err := c.Col().Find(ctx, f)
	errx.CheckError(err)
	data := make([]TEntity, 0)
	err = result.All(context.Background(), &data)
	errx.CheckError(err)
	return data
}

func (c *Collection[TEntity]) GetList(ctx context.Context, filter bson.D, p *x.PageParam, opt *options.FindOptions) ([]TEntity, int64) {
	totalCount := c.Count(ctx, filter)
	filter = c.MergeGlobalFilter(ctx, filter)
	findOptions := options.Find().SetLimit(p.Limit()).SetSkip(p.Skip())
	data := c.Find(ctx, filter, findOptions, opt)
	return data, totalCount
}

func (c *Collection[TEntity]) Exists(ctx context.Context, id primitive.ObjectID) bool {
	errx.NotNil(id, "id")
	filter := c.setSoftDeleteFilter(bson.D{{Key: "_id", Value: id}})
	filter = c.MergeGlobalFilter(ctx, filter)
	var result TEntity
	err := c.Col().FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		panic(err)
	}
	return true
}

func (c *Collection[TEntity]) Insert(ctx context.Context, entity *TEntity) primitive.ObjectID {
	ddd.SetAudited(ctx, entity)
	result, err := c.Col().InsertOne(ctx, entity)
	errx.CheckError(err)
	return result.InsertedID.(primitive.ObjectID)
}

func (c *Collection[TEntity]) InsertMany(ctx context.Context, entitis []TEntity, ignoreErr bool) []primitive.ObjectID {
	if len(entitis) == 0 {
		return make([]primitive.ObjectID, 0)
	}
	data := ddd.SetAuditedMany(ctx, entitis)
	opts := options.InsertMany().SetOrdered(!ignoreErr)
	result, err := c.Col().InsertMany(ctx, data, opts)
	if !ignoreErr {
		errx.CheckError(err)
	} else {
		log.Warn(err.Error())
	}
	if result == nil {
		return make([]primitive.ObjectID, 0)
	}
	return linq.Map(result.InsertedIDs, func(t *interface{}) primitive.ObjectID { return (*t).(primitive.ObjectID) })
}

func (c *Collection[TEntity]) UpdateByID(ctx context.Context, id primitive.ObjectID, entity *TEntity) int {
	ddd.SetAudited(ctx, entity)
	result, err := c.Col().UpdateByID(ctx, id, bson.D{{Key: "$set", Value: entity}})
	errx.CheckError(err)
	return int(result.ModifiedCount)
}

func (c *Collection[TEntity]) UpdateOne(ctx context.Context, filter bson.D, entity *TEntity, opts ...*options.UpdateOptions) int {
	filter = c.setSoftDeleteFilter(filter)
	ddd.SetAudited(ctx, entity)
	result, err := c.Col().UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: entity}}, opts...)
	errx.CheckError(err)
	return int(result.ModifiedCount)
}

func (c *Collection[TEntity]) Update(ctx context.Context, filter bson.D, update interface{}, opts ...*options.UpdateOptions) int {
	filter = c.setSoftDeleteFilter(filter)
	result, err := c.Col().UpdateMany(ctx, filter, update, opts...)
	errx.CheckError(err)
	return int(result.ModifiedCount)
}

func (c *Collection[TEntity]) Delete(ctx context.Context, id primitive.ObjectID) int {
	var v any = (*TEntity)(nil)
	if _, ok := v.(ddd.SoftDeleteEntity); ok {
		e := c.Get(ctx, id)
		var v any = e
		softEntity := v.(ddd.DeletionAuditedEntity)
		softEntity.Deleting(ctx)
		return int(c.UpdateByID(ctx, id, e))
	} else {
		result, err := c.Col().DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
		errx.CheckError(err)
		return int(result.DeletedCount)
	}
}

func (c *Collection[TEntity]) DeleteMany(ctx context.Context, ids []primitive.ObjectID) int {
	var v any = (*TEntity)(nil)
	if _, ok := v.(ddd.DeletionAuditedEntity); ok {
		ctx := ginx.WithScopedUnitwork(ctx)
		var count int = 0
		for i := range ids {
			e := c.Get(ctx, ids[i])
			var v any = e
			softEntity := v.(ddd.DeletionAuditedEntity)
			softEntity.Deleting(ctx)
			count += c.UpdateByID(ctx, ids[i], e)
		}
		return int(count)
	} else {
		result, err := c.Col().DeleteMany(ctx, bson.D{{Key: "_id", Value: bson.M{"$in": ids}}})
		errx.CheckError(err)
		return int(result.DeletedCount)
	}
}

func (c *Collection[TEntity]) DeleteByFilter(ctx context.Context, filter bson.D) int {
	var v any = (*TEntity)(nil)
	if _, ok := v.(ddd.DeletionAuditedEntity); ok {
		es := c.Find(ctx, filter, nil)
		ctx := ginx.WithScopedUnitwork(ctx)
		var count int = 0
		for i := range es {
			var v any = &es[i]
			softEntity := v.(ddd.DeletionAuditedEntity)
			softEntity.Deleting(ctx)
			count += c.UpdateByID(ctx, softEntity.GetId(), &es[i])
		}
		return count
	} else {
		result, err := c.Col().DeleteMany(ctx, filter)
		errx.CheckError(err)
		return int(result.DeletedCount)
	}
}

func (c *Collection[TEntity]) setSoftDeleteFilter(filter bson.D) bson.D {
	var v any = (*TEntity)(nil)
	if _, ok := v.(ddd.SoftDeleteEntity); ok {
		return append(filter, bson.E{Key: ddd.SoftDeleteFieldName, Value: false})
	}
	return filter
}

func (c *Collection[TEntity]) Col() *mongo.Collection {
	return c.c
}

func (c *Collection[TEntity]) MergeGlobalFilter(ctx context.Context, filter bson.D) bson.D {
	dfs := store.DataFilters[TEntity](ctx, c.opts)
	for _, v := range dfs {
		df := v.Filter(ctx, filter)
		v, ok := df.(bson.D)
		if !ok {
			panic(errors.New("data filter type error: must be bson.D"))
		}
		filter = v
	}
	return filter
}

func MergeGlobalFilter[T any](ctx context.Context, filter bson.D, opts *store.StoreOptions) bson.D {
	dfs := store.DataFilters[T](ctx, opts)
	for _, v := range dfs {
		df := v.Filter(ctx, filter)
		if v, ok := df.(bson.D); ok {
			filter = append(filter, v...)
		}
	}
	return filter
}
