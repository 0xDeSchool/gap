package ginx

import (
	"context"
	"errors"
	"net/http"

	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type ErrorHandFunc func(*gin.Context, ...error)

const ErrHandlersKey = "ErrorHandlers"

type ErrorHandlers struct {
	handlers []ErrorHandFunc
}

func NewErrorHandlers() *ErrorHandlers {
	return &ErrorHandlers{
		handlers: make([]ErrorHandFunc, 0),
	}
}

func (eh *ErrorHandlers) Run(c *gin.Context, errs ...error) {
	for i := range eh.handlers {
		eh.handlers[i](c, errs...)
	}
}
func (eh *ErrorHandlers) Add(h ErrorHandFunc) {
	eh.handlers = append(eh.handlers, h)
}

// ErrorMiddleware request panic error handler
func ErrorMiddleware(c *gin.Context) {
	handlers := &ErrorHandlers{}
	c.Set(ErrHandlersKey, handlers)
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				if s, ok := r.(string); ok {
					err = errors.New(s)
				} else {
					err = errors.New("未知错误")
				}
			}
			log.Error("请求出现错误", err)
			handlers.Run(c, err)
			if !c.IsAborted() {
				if err2, ok := r.(*errx.HttpError); ok {
					JSON(c, err2.HttpStatus, *err2)
				} else {
					if errors.Is(err, mongo.ErrNoDocuments) {
						JSON(c, http.StatusNotFound, errx.New(err))
					} else {
						JSON(c, http.StatusInternalServerError, errx.New(err))
					}
				}
				c.Abort()
			}
		}
	}()
	c.Next()
}

// OnError 添加请求生命周期中出现错误时的处理方法。
// 注意: 该方法中不要使用panic
func OnError(ctx context.Context, h ErrorHandFunc) {
	v := ctx.Value(ErrHandlersKey)
	if v != nil {
		v.(*ErrorHandlers).Add(h)
	}
}

func Error(ctx *gin.Context, err *errx.HttpError) {
	JSON(ctx, err.HttpStatus, err)
}

func NotFound(ctx *gin.Context) {
	Error(ctx, errx.ErrPageNotFound)
}

func EntityNotFound(ctx *gin.Context, msg string) {
	Error(ctx, &errx.HttpError{
		HttpStatus: http.StatusNotFound,
		Code:       errx.ErrCodePageNotFound,
		Message:    msg,
	})
}
