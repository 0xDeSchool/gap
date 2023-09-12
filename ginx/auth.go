package ginx

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/errx"
	"github.com/gin-gonic/gin"
)

type AuthHandlerContext[TKey comparable] struct {
	ctx        *gin.Context
	User       *CurrentUserInfo[TKey]
	HasHandled bool
}

func (ac *AuthHandlerContext[TKey]) Context() *gin.Context {
	return ac.ctx
}

type AuthHandler[TKey comparable] func(ctx *AuthHandlerContext[TKey]) error

func AuthFunc[TKey comparable]() gin.HandlerFunc {
	return func(c *gin.Context) {
		handlers := app.GetArray[AuthHandler[TKey]]()
		ctx := &AuthHandlerContext[TKey]{
			ctx: c,
		}
		for i := range handlers {
			err := handlers[i](ctx)
			errx.CheckError(err)
			if ctx.HasHandled {
				break
			}
		}
		if ctx.User == nil {
			ctx.User = &CurrentUserInfo[TKey]{}
		}
		userCtx := context.WithValue(c.Request.Context(), userKey, ctx.User)
		c.Request = c.Request.WithContext(userCtx)
		c.Next()
	}
}
