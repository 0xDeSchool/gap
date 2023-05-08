package ginx

import (
	"github.com/0xDeSchool/gap/app"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type DataResult[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

func JSONResult[T any](ctx *gin.Context, v T) {
	JSON(ctx, http.StatusOK, DataResult[T]{
		Data:    v,
		Success: true,
	})
}

type ResponseResult struct {
	Value any `json:"value"`
}

func NewResponseResult(v any) *ResponseResult {
	if v != nil {
		value := reflect.ValueOf(v)
		if value.Kind() == reflect.Struct {
			ptr := reflect.New(value.Type())
			ptr.Elem().Set(value)
			v = ptr.Interface()
		}
	}
	return &ResponseResult{Value: v}
}

type ResultHandler func(ctx *gin.Context, v *ResponseResult)

func JSON(ctx *gin.Context, code int, v any) {
	handlers := app.GetArray[*ResultHandler]()
	res := NewResponseResult(v)
	for _, h := range handlers {
		(*h)(ctx, res)
	}
	ctx.JSON(code, res.Value)
}

func Ok(ctx *gin.Context) {
	JSON(ctx, http.StatusOK, DataResult[struct{}]{
		Success: true,
		Data:    struct{}{},
	})
}

type CreateEntityResult struct {
	ID string `json:"id"`
}

type CreateManyEntitiesResult struct {
	Count int `json:"count"`
}

type EntityUpdatedResult struct {
	Count int `json:"count"`
}

type SuccessWithMessageResult struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

func EntityCreated(ctx *gin.Context, id string) {
	ctx.JSON(http.StatusOK, CreateEntityResult{
		ID: id,
	})
}

func ManyEntitiesCreated(ctx *gin.Context, count int) {
	ctx.JSON(http.StatusOK, CreateManyEntitiesResult{
		Count: count,
	})
}

func EntityUpdated(ctx *gin.Context, count int) {
	ctx.JSON(http.StatusOK, EntityUpdatedResult{
		Count: count,
	})
}

func SuccessWithMessage(ctx *gin.Context, success bool, msg string) {
	ctx.JSON(http.StatusOK, SuccessWithMessageResult{
		Success: success,
		Msg:     msg,
	})
}
