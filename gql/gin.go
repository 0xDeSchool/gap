package gql

import (
	"context"
	"github.com/gin-gonic/gin"
)

type ginContextKeyType string

const ginContextKey ginContextKeyType = "ginContext"

func GinContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newCtx := context.WithValue(ctx.Request.Context(), ginContextKey, ctx)
		ctx.Request = ctx.Request.WithContext(newCtx)
		ctx.Next()
	}
}

func GinContext(ctx context.Context) *gin.Context {
	return ctx.Value(ginContextKey).(*gin.Context)
}
