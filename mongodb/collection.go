package mongodb

import (
	"context"
	"errors"
	"github.com/0xDeSchool/gap/multi_tenancy"
	"strings"

	"github.com/0xDeSchool/gap/ginx"
	"github.com/0xDeSchool/gap/log"
	"github.com/0xDeSchool/gap/x"

	"github.com/0xDeSchool/gap/ddd"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[TEntity any, TKey comparable] struct {
	c    *mongo.Collection
	opts *store.StoreOptions
}

func NewCollection[TEntity any, TKey comparable](c *mongo.Collection, opts *store.StoreOptions) *Collection[TEntity, TKey] {
	return &Collection[TEntity, TKey]{
		c:    c,
		opts: opts,
	}
}

func (c *Collection[TEntity, TKey]) Find(ctx context.Context, filter bson.D, opts ...*options.FindOptions) ([]TEntity, error) {
	filter = c.SetAllFilter(ctx, filter)
	cur, err := c.Col().Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	data := make([]TEntity, 0)
	err = cur.All(context.Background(), &data)
	return data, err
}

func (c *Collection[TEntity, TKey]) FindByPage(ctx context.Context, filter bson.D, p *x.PageAndSort, opts ...*options.FindOptions) (*x.PagedResult[TEntity], error) {
	result := &x.PagedResult[TEntity]{}
	findOptions := options.Find()
	if p != nil {
		findOptions.SetLimit(p.Limit() + 1).SetSkip(p.Skip())
		if p.IncludeTotal {
			total, err := c.Count(ctx, filter)
			if err != nil {
				return nil, err
			}
			result.Total = total
		}
		sort := c.ParseSort(p)
		findOptions.SetSort(sort)
	}
	newOpts := []*options.FindOptions{findOptions}
	newOpts = append(newOpts, opts...)
	data, err := c.Find(ctx, filter, newOpts...)
	if err != nil {
		return nil, err
	}
	if p != nil && len(data) > int(p.Limit()) {
		data = data[:p.Limit()]
		result.HasMore = true
	}
	return result, nil
}

func (c *Collection[TEntity, TKey]) FindOne(ctx context.Context, filter bson.D, opts ...*options.FindOneOptions) (*TEntity, error) {
	filter = c.SetAllFilter(ctx, filter)
	result := c.Col().FindOne(ctx, filter, opts...)
	err := result.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errx.DataNotFoundError
	} else if err != nil {
		return nil, err
	}
	var v TEntity
	err = result.Decode(&v)
	return &v, err
}

func (c *Collection[TEntity, TKey]) GetMany(ctx context.Context, ids []string) ([]TEntity, error) {
	if len(ids) == 0 {
		return make([]TEntity, 0), nil
	}
	f := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	return c.Find(ctx, f)
}

func (c *Collection[TEntity, TKey]) Count(ctx context.Context, filter bson.D, opts ...*options.CountOptions) (int64, error) {
	filter = c.SetAllFilter(ctx, filter)
	totalCount, err := c.Col().CountDocuments(ctx, filter, opts...)
	return totalCount, err
}

func (c *Collection[TEntity, TKey]) Insert(ctx context.Context, entity *TEntity, opts ...*options.InsertOneOptions) (*TEntity, error) {
	ddd.SetAudited(ctx, entity)
	_, err := c.Col().InsertOne(ctx, entity, opts...)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (c *Collection[TEntity, TKey]) InsertMany(ctx context.Context, entities []TEntity, ignoreErr bool, opts ...*options.InsertManyOptions) ([]TEntity, error) {
	if len(entities) == 0 {
		return entities, nil
	}
	data := ddd.SetAuditedMany(ctx, entities)
	opt := options.InsertMany().SetOrdered(!ignoreErr)
	opts = append(opts, opt)
	result, err := c.Col().InsertMany(ctx, data, opts...)
	if err != nil {
		if !ignoreErr {
			return nil, err
		} else {
			log.Warn("ignored mongodb insert many error: " + err.Error())
		}
	}
	if result == nil {
		return entities, nil
	}
	return entities, nil
}

func (c *Collection[TEntity, TKey]) UpdateOne(ctx context.Context, filter bson.D, entity *TEntity, opts ...*options.UpdateOptions) (int, error) {
	filter = c.SetAllFilter(ctx, filter)
	ddd.SetAudited(ctx, entity)
	result, err := c.Col().UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: entity}}, opts...)
	if err != nil {
		return 0, err
	}
	return int(result.ModifiedCount), nil
}

func (c *Collection[TEntity, TKey]) UpdateMany(ctx context.Context, filter bson.D, update interface{}, opts ...*options.UpdateOptions) (int, error) {
	filter = c.SetAllFilter(ctx, filter)
	set := bson.D{{Key: "$set", Value: update}}
	result, err := c.Col().UpdateMany(ctx, filter, set, opts...)
	if err != nil {
		return 0, err
	}
	return int(result.ModifiedCount), nil
}

func (c *Collection[TEntity, TKey]) DeleteOne(ctx context.Context, filter bson.D) (int, error) {
	filter = c.SetAllFilter(ctx, filter)
	var v any = (*TEntity)(nil)
	if _, ok := v.(ddd.ISoftDeleteEntity); ok {
		e, err := c.FindOne(ctx, filter)
		if err != nil {
			return 0, err
		}
		var se any = e
		softEntity := se.(ddd.ISoftDeleteEntity)
		softEntity.Deleting(ctx)
		count, err := c.UpdateOne(ctx, filter, e)
		if err != nil {
			return 0, err
		}
		return count, nil
	} else {
		result, err := c.Col().DeleteOne(ctx, filter)
		if err != nil {
			return 0, err
		}
		return int(result.DeletedCount), nil
	}
}

func (c *Collection[TEntity, TKey]) DeleteMany(ctx context.Context, filter bson.D) (int, error) {
	filter = c.SetAllFilter(ctx, filter)
	var v any = (*TEntity)(nil)
	if _, ok := v.(ddd.ISoftDeleteEntity); ok {
		es, err := c.Find(ctx, filter, nil)
		if err != nil {
			return 0, err
		}
		uw := ginx.NewUnitWork(ctx)
		uw.Start(ctx)
		var count = 0
		for i := range es {
			var v any = &es[i]
			softEntity := v.(ddd.ISoftDeleteEntity)
			softEntity.Deleting(ctx)
			idFilter := bson.D{{"_id", softEntity.GetId()}}
			ct, err := c.UpdateOne(ctx, idFilter, &es[i])
			if err != nil {
				uw.Abort(ctx)
				return 0, err
			}
			count += ct
		}
		err = uw.Commit(ctx)
		if err != nil {
			return 0, err
		}
		return count, nil
	} else {
		result, err := c.Col().DeleteMany(ctx, filter)
		if err != nil {
			return 0, nil
		}
		return int(result.DeletedCount), nil
	}
}

func (c *Collection[TEntity, TKey]) Col() *mongo.Collection {
	return c.c
}

func (c *Collection[TEntity, TKey]) SetAllFilter(ctx context.Context, filter bson.D) bson.D {
	filter = c.MergeGlobalFilter(ctx, filter)
	filter = c.SetSoftDeleteFilter(filter)
	filter = c.SetMultiTenantFilter(ctx, filter)
	return filter
}

func (c *Collection[TEntity, TKey]) SetSoftDeleteFilter(filter bson.D) bson.D {
	var v any = (*TEntity)(nil)
	if _, ok := v.(ddd.SoftDeleteEntity); ok {
		return append(filter, bson.E{Key: ddd.SoftDeleteFieldName, Value: false})
	}
	return filter
}

func (c *Collection[TEntity, TKey]) SetMultiTenantFilter(ctx context.Context, filter bson.D) bson.D {
	enabled := multi_tenancy.IsEnableMultiTenantFilter(ctx)
	if !enabled {
		return filter
	}
	var v any = (*TEntity)(nil)
	if _, ok := v.(ddd.IMultiTenancy); ok {
		return append(filter, bson.E{Key: ddd.TenantIdDbKey, Value: multi_tenancy.CurrentTenant(ctx).Id})
	}
	return filter
}

func (c *Collection[TEntity, TKey]) MergeGlobalFilter(ctx context.Context, filter bson.D) bson.D {
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

func (c *Collection[TEntity, TKey]) ParseSort(p *x.PageAndSort) bson.D {
	sort := bson.D{}
	if p.Sort != "" {
		desc := 1
		k := ""
		if p.IsDesc() {
			desc = -1
			k = p.Sort[1:]
		} else {
			k = strings.TrimLeft(p.Sort, "+")
		}
		sort = append(sort, bson.E{Key: k, Value: desc})
	}
	if x.CanConvert[TEntity, ddd.CreationAuditedEntity]() && !strings.Contains(p.Sort, "createdAt") {
		sort = append(sort, bson.E{Key: "createdAt", Value: -1})
	}
	return sort
}

func MergeFilter[T any](ctx context.Context, filter bson.D, opts *store.StoreOptions) bson.D {
	dfs := store.DataFilters[T](ctx, opts)
	for _, v := range dfs {
		df := v.Filter(ctx, filter)
		if v, ok := df.(bson.D); ok {
			filter = append(filter, v...)
		}
	}
	return filter
}
