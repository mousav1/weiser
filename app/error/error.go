package error

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

type Error struct {
	Message string
}
type wrappedError struct {
	cause error
	msg   string
	stack []uintptr
	code  int
}

func (e *Error) Error() string {
	return e.Message
}

func New(message string) error {
	return &Error{Message: message}
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Newf(format string, args ...interface{}) error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

func (e *wrappedError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.cause)
}

func (e *wrappedError) Unwrap() error {
	return e.cause
}

func (e *wrappedError) Code() int {
	return e.code
}

func (e *wrappedError) WithCode(code int) error {
	return &wrappedError{
		cause: e.cause,
		msg:   e.msg,
		stack: e.stack,
		code:  code,
	}
}

func (e *wrappedError) WithCause(cause error) error {
	return &wrappedError{
		cause: cause,
		msg:   e.msg,
		stack: e.stack,
		code:  e.code,
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	return &wrappedError{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func GetRootCause(err error) error {
	for {
		if u := errors.Unwrap(err); u != nil {
			err = u
		} else {
			break
		}
	}
	return err
}

func (e *Error) WithMessage(message string) error {
	return &Error{Message: message + ": " + e.Error()}
}

func (e *Error) WithStack() error {
	return errors.WithStack(e)
}

func (e *Error) Unwrap() error {
	return errors.Cause(e)
}
