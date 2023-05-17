package mongodb

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/ginx"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUnitWork struct {
	session mongo.Session
}

// Abort implements UnitWork
func (uw *mongoUnitWork) Abort(ctx context.Context) {
	uw.session.AbortTransaction(ctx)
	uw.session.EndSession(ctx)
}

// Commit implements UnitWork
func (uw *mongoUnitWork) Commit(ctx context.Context) error {
	err := uw.session.CommitTransaction(ctx)
	uw.session.EndSession(ctx)
	return err
}

// Start implements UnitWork
func (uw *mongoUnitWork) Start(ctx context.Context) context.Context {
	err := uw.session.StartTransaction()
	errx.CheckError(err)
	return mongo.NewSessionContext(ctx, uw.session)
}

func NewMongoUnitWork() *ginx.UnitWork {
	client := app.Get[mongo.Client]()
	session, err := client.StartSession()
	errx.CheckError(err)
	var uw ginx.UnitWork = &mongoUnitWork{
		session: session,
	}
	return &uw
}
