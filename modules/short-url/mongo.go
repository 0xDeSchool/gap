package short_url

import (
	"context"
	"github.com/0xDeSchool/gap/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

const SHortUrlCollection = "short_urls"

type mongoRepository struct {
	*mongodb.MongoRepositoryBase[ShortUrl, string]
}

func NewMongoRepository() *ShortUrlRepository {
	var r ShortUrlRepository = &mongoRepository{
		MongoRepositoryBase: mongodb.NewMongoRepositoryBase[ShortUrl, string](SHortUrlCollection),
	}
	return &r
}

func (r *mongoRepository) GetUrl(ctx context.Context, key string) (*ShortUrl, error) {
	filter := bson.D{{"key", key}}
	urls, err := r.Collection(ctx).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(urls) == 0 {
		return nil, nil
	} else {
		return urls[0], nil
	}
}
