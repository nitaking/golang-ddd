package application_error

import "fmt"

type Code string

const (
	CodeInvalidRequest Code = "invalid_request"
	CodeNotFound       Code = "not_found"
	CodeConflict       Code = "conflict"
	CodeUnexpected     Code = "unexpected"
)

type AppError struct {
	code    Code
	message string
	cause   error
}

func New(code Code, message string) *AppError {
	return &AppError{code: code, message: message}
}

func (e *AppError) Wrap(err error) *AppError {
	return &AppError{code: e.code, message: e.message, cause: err}
}

func (e *AppError) Error() string {
	if e.cause == nil {
		return e.message
	}
	return fmt.Sprintf("%s: %v", e.message, e.cause)
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func (e *AppError) Code() Code {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}
