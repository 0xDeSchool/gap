package mongodb

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/store"
	"github.com/0xDeSchool/gap/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	app.Configure(func() error {
		var mongoOptions = &MongoOptions{
			URL: "mongodb://localhost:27017",
		}
		utils.ViperBind("Mongo", mongoOptions)
		app.TryAddValue(mongoOptions)

		stp := store.NewStoreOptions()
		app.TryAddValue(stp)

		app.TryAddSingleton(func() *mongo.Client {
			c, err := GetClient(mongoOptions)
			errx.CheckError(err)
			return c
		})
		app.TryAddTransient(NewMongoUnitWork)
		return nil
	})
	app.PostRun(func() error {
		c := app.Get[mongo.Client]()
		return c.Disconnect(context.Background())
	})
}
