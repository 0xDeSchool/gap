package gql

import (
	"context"
	"errors"
	"fmt"
	"github.com/0xDeSchool/gap/ginx"
	"github.com/0xDeSchool/gap/log"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AddGraphQL(s *ginx.Server, gqlServer *handler.Server) {
	s.G.Use(GinContextMiddleware())
	// jwt middleware
	s.G.POST("/query", func(ctx *gin.Context) {
		gqlServer.ServeHTTP(ctx.Writer, ctx.Request)
	})
	h := playground.Handler("GraphQL playground", "/query")
	s.G.GET("/", func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	})
	gqlServer.SetRecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
		msg := fmt.Sprintf("%s", err)
		log.Error(msg)
		if e, ok := err.(error); ok {
			return gqlerror.Errorf(e.Error())
		}
		return gqlerror.Errorf(msg)
	})
}

func RequireAuthDirective[TID comparable](ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	currentUser := ginx.CurrentUser[TID](ctx)
	if !currentUser.Authenticated() {
		return nil, errors.New("not authenticated")
	}
	return next(ctx)
}
