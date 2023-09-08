package ginx

import (
	"context"
	"errors"
	"github.com/0xDeSchool/gap/app"

	"github.com/0xDeSchool/gap/errx"
	"github.com/gin-gonic/gin"
)

const UnitWorkKey = "UnitWorkKey"

type UnitWork interface {
	Start(ctx context.Context) context.Context
	Commit(ctx context.Context) error
	Abort(ctx context.Context)
}

type unitWorkManager struct {
	unitworks []UnitWork
}

func newUnitWorkManager() *unitWorkManager {
	return &unitWorkManager{
		unitworks: []UnitWork{nil},
	}
}

func (uw *unitWorkManager) Default() UnitWork {
	if uw.unitworks[0] == nil {
		uw.unitworks[0] = *app.Get[UnitWork]()
	}
	return uw.unitworks[0]
}

func (uw *unitWorkManager) New() UnitWork {
	nu := *app.Get[UnitWork]()
	uw.unitworks = append(uw.unitworks, nu)
	return nu
}

func (uw *unitWorkManager) CommitAll(ctx context.Context) []error {
	errs := make([]error, 0)
	for i := range uw.unitworks {
		if uw.unitworks[i] != nil {
			err := uw.unitworks[i].Commit(ctx)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func (uw *unitWorkManager) AbortAll(ctx context.Context) {
	for i := range uw.unitworks {
		if uw.unitworks[i] != nil {
			uw.unitworks[i].Abort(ctx)
		}
	}
}

func UnitWorkMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		um := newUnitWorkManager()
		newCtx := context.WithValue(ctx.Request.Context(), UnitWorkKey, um)
		ctx.Request = ctx.Request.WithContext(newCtx)
		OnError(ctx, func(ctx *gin.Context, errs ...error) {
			um.AbortAll(context.Background())
		})
		ctx.Next()
		errs := um.CommitAll(context.Background())
		if len(errs) > 0 {
			panic(errx.Errors("commit unitwork error: ", errs...))
		}
	}
}

// WithScopedUnitwork 使用默认的UnitWork，在当前请求周期中有效，自动提交事务
func WithScopedUnitwork(ctx context.Context) context.Context {
	v := ctx.Value(UnitWorkKey)
	if v != nil {
		uwm := v.(*unitWorkManager)
		return uwm.Default().Start(ctx)
	} else {
		panic(errors.New("context has no unitwork instance"))
	}
}

func AbortScopedUnitWork(ctx context.Context) {
	v := ctx.Value(UnitWorkKey)
	if v != nil {
		uwm := v.(*unitWorkManager)
		uwm.Default().Abort(ctx)
	}
}

// NewUnitWork 在当前请求周期中创建新的UnitWork，可控制UnitWork提交或取消
// 另外，在请求结束时，UnitWorkMiddleware会自动提交，当请求执行过程中发生错误时自动abort
func NewUnitWork(ctx context.Context) UnitWork {
	v := ctx.Value(UnitWorkKey)
	if v == nil {
		v = newUnitWorkManager()
		ctx = context.WithValue(ctx, UnitWorkKey, v)
	}
	uwm := v.(*unitWorkManager)
	return uwm.New()
}
