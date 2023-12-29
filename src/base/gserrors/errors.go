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

// gsError 自定义错误对象
type errorHost struct {
	origin  error  // origin error
	stack   string // stack trace message
	message string // error message
}

func (err *errorHost) Error() string {
	if err.message != "" {
		if err.origin != nil {
			return fmt.Sprintf(
				"%s\norigin: %s\nbacktrace:\n%s",
				err.message,
				err.origin.Error(),
				err.stack,
			)
		}
		return fmt.Sprintf("%s\n%s", err.message, err.stack)
	}
	if err.origin != nil {
		return fmt.Sprintf(
			"origin: %s\nbacktrace:\n%s",
			err.origin.Error(),
			err.stack,
		)
	}
	return fmt.Sprintf("<unknown error>\nbacktrace:\n%s", err.stack)
}

func (err *errorHost) String() string {
	return err.Error()
}

func (err *errorHost) Stack() string {
	return err.stack
}

func (err *errorHost) Origin() error {
	return err.origin
}

func (err *errorHost) NewOrigin(origin error) {
	err.origin = origin
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
func New(message string) GSError {
	return &errorHost{
		origin:  nil,
		message: message,
		stack:   string(stack()),
	}
}

// Newf create new GSError object
func Newf(template string, args ...interface{}) GSError {
	return &errorHost{
		origin:  nil,
		stack:   string(stack()),
		message: fmt.Sprintf(template, args...),
	}
}

// Panic create new GSError and panic
func Panic(message string) {
	panic(New(message))
}

// PanicWith create new GSError and panic
func PanicWith(err error, message string) {
	panic(NewWith(err, message))
}

// Panicf create new GSError and panic
func Panicf(template string, args ...interface{}) {
	panic(Newf(template, args...))
}

// PanicfWith create new GSError and panic
func PanicfWith(err error, template string, args ...interface{}) {
	panic(NewfWith(err, template, args...))
}

// NewWith create new GSError with origin error
func NewWith(err error, message string) GSError {
	return &errorHost{
		origin:  err,
		message: message,
		stack:   string(stack()),
	}
}

// NewfWith create new GSError with origin error
func NewfWith(err error, template string, args ...interface{}) GSError {
	return &errorHost{
		origin:  err,
		message: fmt.Sprintf(template, args...),
		stack:   string(stack()),
	}
}
