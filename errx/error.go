package errx

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/0xDeSchool/gap/log"
)

type ErrCode = string

var (

	// ErrCodeUnknown 未知错误
	ErrCodeUnknown ErrCode = "unkown"

	// ErrUnknownParameterInvalid 未知原因参数不合法
	ErrUnknownParameterInvalid ErrCode = "invalid"
)

var (
	DataNotFoundError = errors.New("data not found")
)

type errorInfo struct {
	message string
	err     error
}

func CatchPanic(message string) {
	if r := recover(); r != nil {
		log.Warn(fmt.Sprintf(message+": %s", r))
	}
}

func PrintStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}
