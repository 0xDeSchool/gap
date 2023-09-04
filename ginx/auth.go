package ginx

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/gin-gonic/gin"
)

type AuthHandlerContext[TKey comparable] struct {
	ctx        *gin.Context
	User       *CurrentUserInfo[TKey]
	HasHandled bool
}

func (ac *AuthHandlerContext[TKey]) Ctx() context.Context {
	return ac.ctx
}

type AuthHandler[TKey comparable] interface {
	Handle(ctx *AuthHandlerContext[TKey])
}

func AuthFunc[TKey comparable]() gin.HandlerFunc {
	return func(c *gin.Context) {
		handlers := app.GetArray[AuthHandler[TKey]]()
		ctx := &AuthHandlerContext[TKey]{
			ctx: c,
		}
		for i := range handlers {
			handlers[i].Handle(ctx)
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
