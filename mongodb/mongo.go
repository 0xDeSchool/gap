package mongodb

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/store"
	"github.com/0xDeSchool/gap/x"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTransactionContext struct {
	ctx mongo.SessionContext
}

func NewMongoTransactionContext(ctx mongo.SessionContext) MongoTransactionContext {
	return MongoTransactionContext{
		ctx: ctx,
	}
}

type MongoRepositoryBase[TEntity any, TKey comparable] struct {
	Options        *MongoOptions
	CollectionName string
	StoreOpts      *store.StoreOptions
}

func NewMongoRepositoryBase[TEntity any, TKey comparable](collectionName string) *MongoRepositoryBase[TEntity, TKey] {
	return &MongoRepositoryBase[TEntity, TKey]{
		Options:        app.Get[MongoOptions](),
		StoreOpts:      app.Get[store.StoreOptions](),
		CollectionName: collectionName,
	}
}

func (mr *MongoRepositoryBase[TEntity, TKey]) GetCollection(ctx context.Context, name string) *mongo.Collection {
	var mc *mongo.Client
	if sessionCtx, ok := ctx.(mongo.SessionContext); ok {
		mc = sessionCtx.Client()
	} else {
		mc = app.Get[mongo.Client]()
	}
	return mc.Database(mr.Options.DbName).Collection(name)
}

func (mr *MongoRepositoryBase[TEntity, TKey]) Collection(ctx context.Context) *Collection[TEntity, TKey] {
	c := mr.GetCollection(ctx, mr.CollectionName)
	return NewCollection[TEntity, TKey](c, mr.StoreOpts)
}

func (mr *MongoRepositoryBase[TEntity, TKey]) GetPagedList(ctx context.Context, p *x.PageAndSort) (*x.PagedResult[TEntity], error) {
	filter := bson.D{}
	return mr.Collection(ctx).FindByPage(ctx, filter, p)
}

func (mr *MongoRepositoryBase[TEntity, TKey]) GetById(ctx context.Context, id TKey) (*TEntity, error) {
	return mr.Collection(ctx).FindOne(ctx, bson.D{{Key: "_id", Value: id}})
}

func (mr *MongoRepositoryBase[TEntity, TKey]) GetOrNilById(ctx context.Context, id TKey) (*TEntity, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	data, err := mr.Collection(ctx).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return data[0], nil
}

func (mr *MongoRepositoryBase[TEntity, TKey]) GetMany(ctx context.Context, ids []TKey) ([]*TEntity, error) {
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	return mr.Collection(ctx).Find(ctx, filter)
}
func (mr *MongoRepositoryBase[TEntity, TKey]) Exists(ctx context.Context, id TKey) (bool, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	data, err := mr.Collection(ctx).Find(ctx, filter)
	if err != nil {
		return false, err
	}
	return len(data) > 0, nil
}

func (mr *MongoRepositoryBase[TEntity, TKey]) Count(ctx context.Context) (int64, error) {
	filter := bson.D{}
	return mr.Collection(ctx).Count(ctx, filter)
}

func (mr *MongoRepositoryBase[TEntity, TKey]) FindByRegex(ctx context.Context, field, regex string, p *x.PageAndSort) (*x.PagedResult[TEntity], error) {
	filter := bson.D{{Key: field, Value: primitive.Regex{Pattern: regex, Options: "i"}}}
	return mr.Collection(ctx).FindByPage(ctx, filter, p)
}

func (mr *MongoRepositoryBase[TEntity, TKey]) Insert(ctx context.Context, entity *TEntity) (*TEntity, error) {
	return mr.Collection(ctx).Insert(ctx, entity)
}

// InsertMany ignoreErr 是否忽略批量插入时的错误, 一般为false, 当导入时忽略重复key的时候可以设为true
func (mr *MongoRepositoryBase[TEntity, TKey]) InsertMany(ctx context.Context, entities []*TEntity, ignoreErr bool) ([]*TEntity, error) {
	return mr.Collection(ctx).InsertMany(ctx, entities, ignoreErr)
}

func (mr *MongoRepositoryBase[TEntity, TKey]) UpdateById(ctx context.Context, id TKey, data *TEntity) (int, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	return mr.Collection(ctx).UpdateOne(ctx, filter, data)
}

func (mr *MongoRepositoryBase[TEntity, TKey]) Delete(ctx context.Context, id TKey) (int, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	return mr.Collection(ctx).DeleteOne(ctx, filter)
}
func (mr *MongoRepositoryBase[TEntity, TKey]) DeleteMany(ctx context.Context, ids []TKey) (int, error) {
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	return mr.Collection(ctx).DeleteMany(ctx, filter)
}

type orderResult struct {
	MaxOrder float64 `bson:"maxOrder"`
}

func (mr *MongoRepositoryBase[TEntity, TKey]) MaxOrder(ctx context.Context, field string, v any) float64 {
	match := bson.D{{Key: "$match", Value: bson.D{{Key: field, Value: v}}}}
	var groupKey any = primitive.NilObjectID
	if field != "" {
		groupKey = "$" + field
	}
	groupMax := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: groupKey},
		{Key: "maxOrder", Value: bson.M{"$max": "$order"}},
	}}}
	aggregate := bson.A{}
	if field != "" {
		aggregate = append(aggregate, match)
	}
	aggregate = append(aggregate, groupMax)
	result, err := mr.Collection(ctx).Col().Aggregate(ctx, aggregate)
	errx.CheckError(err)
	results := make([]orderResult, 0)
	errx.CheckError(result.All(ctx, &results))
	if len(results) > 0 {
		return results[0].MaxOrder
	}
	return 0
}

func (mr *MongoRepositoryBase[TEntity, TKey]) MaxOrderMany(ctx context.Context, field string, v any) map[any]float64 {
	match := bson.D{{Key: "$match", Value: bson.D{{Key: field, Value: bson.M{"$in": v}}}}}
	var groupKey any = primitive.NilObjectID
	if field != "" {
		groupKey = "$" + field
	}
	groupMax := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: groupKey},
		{Key: "maxOrder", Value: bson.M{"$max": "$order"}},
	}}}
	aggregate := bson.A{}
	if field != "" {
		aggregate = append(aggregate, match)
	}
	aggregate = append(aggregate, groupMax)
	result, err := mr.Collection(ctx).Col().Aggregate(ctx, aggregate)
	errx.CheckError(err)
	var results []bson.M
	errx.CheckError(result.All(ctx, &results))
	data := make(map[any]float64)
	for k := range results {
		key := results[k]["_id"]
		if key != nil {
			if maxOrder, ok := results[k]["maxOrder"].(float64); ok {
				data[key] = maxOrder
			} else if mo, ok := results[k]["maxOrder"].(int32); ok {
				data[key] = float64(mo)
			}
		}
	}
	return data
}
