package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	code    Code
	message string
	cause   error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.cause == nil {
		return e.message
	}
	return fmt.Sprintf("%s: %v", e.message, e.cause)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *Error) Code() Code {
	if e == nil {
		return CodeUnknown
	}
	return e.code
}

func (e *Error) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

func New(code Code, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

func NewValidationError(message string) *Error {
	return New(CodeValidation, message)
}

func NewDomainError(message string) *Error {
	return New(CodeDomainFailure, message)
}

func NewInternalError(message string) *Error {
	return New(CodeInternal, message)
}

func Wrap(err error, code Code, message string) *Error {
	if err == nil {
		return &Error{
			code:    code,
			message: message,
		}
	}
	return &Error{
		code:    code,
		message: message,
		cause:   err,
	}
}

func CodeOf(err error) Code {
	if err == nil {
		return CodeUnknown
	}
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code()
	}
	return CodeUnknown
}

func IsCode(err error, code Code) bool {
	return CodeOf(err) == code
}

func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
