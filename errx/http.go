package errx

import (
	"fmt"
	"net/http"
)

var (
	ErrCodePageNotFound   ErrCode = "Server.PageNotFound"
	ErrCodeEntityNotFound ErrCode = "Server.EntityNotFound"

	// ErrCodeBadRequest 参数校验失败
	ErrCodeBadRequest ErrCode = "Server.BadRequest"
	// ErrCodeUnAuthenticated 未认证(登录)
	ErrCodeUnAuthenticated ErrCode = "Server.UnAuthenticated"
	ErrCodeUnauthorized    ErrCode = "Server.Unauthorized"
)

type HttpError struct {
	HttpStatus int     `json:"-"`
	Code       ErrCode `json:"code"`
	Message    string  `json:"message"`
	Data       any     `json:"data"`
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("http error: code: %s, msg: %s, httpstatus:%d", e.Code, e.Message, e.HttpStatus)
}

var UnhandledError = &HttpError{
	Message:    "服务器未处理异常",
	Code:       ErrCodeUnkown,
	HttpStatus: http.StatusInternalServerError,
}

var ErrUnAuthenticated = &HttpError{
	Message:    "user not login",
	Code:       ErrCodeUnAuthenticated,
	HttpStatus: http.StatusUnauthorized,
}

var ErrUnauthorized = &HttpError{
	Message:    "permission forbidden",
	Code:       ErrCodeUnauthorized,
	HttpStatus: http.StatusUnauthorized,
}

func New(err error) HttpError {
	return HttpError{
		Message: err.Error(),
		Code:    ErrCodeUnkown,
	}
}

var ErrPageNotFound = &HttpError{
	HttpStatus: http.StatusNotFound,
	Code:       ErrCodePageNotFound,
	Message:    "page not found",
}

func NewUnkonwErr(err any) HttpError {
	return HttpError{
		Message: fmt.Sprintf("未处理异常: %s", err),
		Code:    ErrCodeUnkown,
	}
}

func HttpPanic(code ErrCode, message string) {
	panic(&HttpError{
		Message:    message,
		Code:       code,
		HttpStatus: http.StatusInternalServerError,
	})
}

func PanicErr(status int, code ErrCode, message string) {
	panic(&HttpError{
		Message:    message,
		Code:       code,
		HttpStatus: status,
	})
}

func PanicNotFound(message string) {
	panic(&HttpError{
		Message:    message,
		Code:       ErrCodePageNotFound,
		HttpStatus: http.StatusNotFound,
	})
}

func PanicEntityNotFound(message string) {
	panic(&HttpError{
		Message:    message,
		Code:       ErrCodeEntityNotFound,
		HttpStatus: http.StatusNotFound,
	})
}

func PanicValidatition(message string) {
	panic(&HttpError{
		Message:    message,
		Code:       ErrCodeBadRequest,
		HttpStatus: http.StatusBadRequest,
	})
}

func PanicUnAuthenticated(message string) {
	panic(ErrUnAuthenticated)
}

func PanicUnAuthorized(message string) {
	panic(&HttpError{
		Message:    message,
		Code:       ErrCodeUnauthorized,
		HttpStatus: http.StatusForbidden,
	})
}
