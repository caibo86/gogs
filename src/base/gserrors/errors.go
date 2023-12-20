// -------------------------------------------
// @file      : errors.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/20 下午6:10
// -------------------------------------------

package gserrors

import (
	"bytes"
	"fmt"
	"runtime"
)

// GSError 自定义错误接口
type GSError interface {
	error            // inherit from system error interface
	Stack() string   // get stack trace message
	Origin() error   // get origin error object
	NewOrigin(error) // reset the origin error
}

type errorHost struct {
	origin  error  // origin error
	stack   string // stack trace message
	message string // error message
}

func (err *errorHost) Error() string {
	if err.message != "" {
		if err.origin != nil {
			return fmt.Sprintf(
				"%s\nbacktrace:\n%sbacktrace error:\n%s",
				err.message,
				err.stack,
				err.origin.Error(),
			)
		}

		return fmt.Sprintf("%s\nbacktrace:\n%s", err.message, err.stack)
	}

	if err.origin != nil {
		return fmt.Sprintf(
			"%s\nbacktrace:\n%s",
			err.origin.Error(),
			err.stack,
		)
	}
	return fmt.Sprintf("<unknown error>\n%s", err.stack)
}

func (err *errorHost) Stack() string {
	return err.stack
}

func (err *errorHost) Origin() error {
	return err.origin
}

func (err *errorHost) NewOrigin(target error) {
	err.origin = target
}

func stack() []byte {
	var buff bytes.Buffer
	for skip := 2; ; skip++ {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		buff.WriteString(fmt.Sprintf("\tfile = %s, line = %d\n", file, line))
	}

	return buff.Bytes()
}

// New create new GSError object
func New(err error) GSError {
	return &errorHost{
		origin: err,
		stack:  string(stack()),
	}
}

// Newf create new GSError object
func Newf(err error, template string, args ...interface{}) GSError {
	return &errorHost{
		origin:  err,
		stack:   string(stack()),
		message: fmt.Sprintf(template, args...),
	}
}
