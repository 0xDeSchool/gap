package short_url

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/ddd"
	"github.com/dineshappavoo/basex"
	"github.com/yitter/idgenerator-go/idgen"
	"math/big"
)

type KeyCreator interface {
	Create(ctx context.Context) string
}

type ShortUrl struct {
	ddd.EntityBase[string] `bson:",inline"`
	Url                    string `bson:"url"`
	Key                    string `bson:"key"`
}

type ShortUrlRepository interface {
	ddd.RepositoryBase[ShortUrl, string]

	GetUrl(ctx context.Context, key string) (*ShortUrl, error)
}

type DefaultKeyCreator struct {
	creator *idgen.DefaultIdGenerator
}

func NewDefaultKeyCreator() *DefaultKeyCreator {
	opts := idgen.NewIdGeneratorOptions(1)
	opts.WorkerIdBitLength = 2
	opts.SeqBitLength = 3
	creator := idgen.NewDefaultIdGenerator(opts)
	var k KeyCreator = &DefaultKeyCreator{
		creator: creator,
	}
	return k.(*DefaultKeyCreator)
}

func (k *DefaultKeyCreator) Create(ctx context.Context) string {
	id := k.creator.NewLong()
	s, _ := basex.EncodeInt(big.NewInt(id))
	return s
}

func CreateUrl(ctx context.Context, urlStr string) (string, error) {
	repo := app.GetPtr[ShortUrlRepository]()
	key := app.GetPtr[KeyCreator]().Create(ctx)
	u := &ShortUrl{
		Url: urlStr,
		Key: key,
	}
	_, err := repo.Insert(ctx, u)
	if err != nil {
		return "", err
	}
	return key, nil
}
