package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/0xDeSchool/gap/utils/linq"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _writer = zerolog.ConsoleWriter{
	Out:        os.Stderr,
	TimeFormat: "2006/01/02 15:04:05",
}
var _logger = log.Output(_writer)

func Writer() io.Writer {
	return _writer
}

func Info(msg string) {
	_logger.Info().Msg(msg)
}

func Infof(format string, v ...interface{}) {
	_logger.Info().Msgf(format, v...)
}

func Debug(msg string) {
	_logger.Debug().Msg(msg)
}

func Debugf(format string, v ...interface{}) {
	_logger.Debug().Msgf(format, v...)
}

func Warn(msg string, errs ...error) {
	lerr := _logger.Warn()
	if len(errs) > 0 {
		//lerr = lerr.Err(errs[0])
		msg += ": " + errs[0].Error()
	}
	lerr.Msg(msg)
}

func Warnf(format string, v ...interface{}) {
	_logger.Warn().Msgf(format, v...)
}

func Error(msg string, errs ...any) {
	errMsg := linq.Map(errs, func(err *any) string { return fmt.Sprintf("%v", *err) })
	msg += msg + ":" + strings.Join(errMsg, ",")
	msg += "\n" + printStack()
	logger := _logger.Error()
	if len(errs) > 0 {
		logger = logger.Err(errors.New(errMsg[0]))
	}
	logger.Msg(msg)
}

func Fatal(err error, msg ...string) {
	msg = append(msg, "\n"+printStack())
	lerr := _logger.Fatal().Err(err)
	if len(msg) > 0 {
		lerr.Msg(strings.Join(msg, ","))
	} else {
		lerr.Msg(err.Error())
	}
}

func FatalMsg(msg string) {
	msg = msg + "\n" + printStack()
	lerr := _logger.Fatal()
	lerr.Msg(msg)
}

func Fatalf(format string, a ...any) {
	lerr := _logger.Fatal()
	lerr.Msg(fmt.Sprintf(format, a...) + "\n" + printStack())
}

func printStack() string {
	return string(debug.Stack())
}
