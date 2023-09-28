package ginx

import (
	"context"
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/eventbus"
	"github.com/gin-gonic/gin"
)

type AuthOptions[TKey comparable] struct {
	handlers []AuthHandler[TKey]
}

type AuthedEventData[TKey comparable] struct {
	ID       TKey           `json:"id"`
	UserName string         `json:"userName"`
	Address  string         `json:"address"`
	Avatar   string         `json:"avatar"`
	Data     map[string]any `json:"data"`
	AuthType string         // 验证类型
}

func NewAuthedEventData[TKey comparable](u *CurrentUserInfo[TKey]) *AuthedEventData[TKey] {
	return &AuthedEventData[TKey]{
		ID:       u.ID,
		UserName: u.UserName,
		Address:  u.Address,
		Avatar:   u.Avatar,
		Data:     u.Data,
		AuthType: u.AuthType,
	}
}

func (o *AuthOptions[TKey]) AddHandler(h AuthHandler[TKey]) {
	o.handlers = append(o.handlers, h)
}

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
		opts := app.Get[AuthOptions[TKey]]()
		handlers := opts.handlers
		ctx := &AuthHandlerContext[TKey]{
			ctx: c,
		}
		for i := range opts.handlers {
			err := handlers[i](ctx)
			errx.CheckError(err)
			if ctx.HasHandled {
				break
			}
		}
		if ctx.User == nil {
			ctx.User = &CurrentUserInfo[TKey]{}
		} else {
			eventbus.Publish(context.Background(), NewAuthedEventData[TKey](ctx.User))
		}

		userCtx := WithUserContext(c.Request.Context(), ctx.User)
		c.Request = c.Request.WithContext(userCtx)
		c.Next()
	}
}

// AddAuthHandler adds a handler to the auth middleware.
func AddAuthHandler[TKey comparable](h AuthHandler[TKey]) {
	app.TryAddValue(&AuthOptions[TKey]{})
	app.ConfigureOptions(func(c *app.Container, opts *AuthOptions[TKey]) {
		opts.AddHandler(h)
	})
}
