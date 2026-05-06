package errx

import (
	"errors"
	"fmt"
)

type Fields map[string]any

type Error struct {
	Code   Code
	Op     string
	Msg    string
	Err    error
	Fields Fields
}

func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s (%s)", e.Op, e.Msg, e.Code)
	}
	return fmt.Sprintf("%s: %s (%s): %v", e.Op, e.Msg, e.Code, e.Err)
}

func (e *Error) Unwrap() error { return e.Err }

func E(op string, code Code, msg string, err error) *Error {
	return &Error{
		Code:   code,
		Op:     op,
		Msg:    msg,
		Err:    err,
		Fields: Fields{},
	}
}

func CodeOf(err error) Code {
	if err == nil {
		return ""
	}
	var ex *Error
	if errors.As(err, &ex) {
		return ex.Code
	}
	return CodeUnknown
}

func IsCode(err error, code Code) bool {
	return CodeOf(err) == code
}
