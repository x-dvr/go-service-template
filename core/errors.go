package core

import "fmt"

// ErrorCode represents type of domain error
type ErrorCode int

const (
	ErrUnknown ErrorCode = iota
	ErrNotFound
	ErrUnprocessable
	ErrValidation
	ErrDuplicate
	ErrUnauthorized
	ErrPermissionDenied
)

// String returns string representation of error code
func (ec ErrorCode) String() string {
	switch ec {
	case ErrNotFound:
		return "not found"
	case ErrUnprocessable:
		return "unprocessable entity"
	case ErrValidation:
		return "validation error"
	case ErrDuplicate:
		return "duplication error"
	case ErrPermissionDenied:
		return "permission denied"
	case ErrUnauthorized:
		return "unauthorized"
	case ErrUnknown:
		fallthrough
	default:
		return "unknown error"
	}
}

// Error represents domain error
type Error struct {
	Code    ErrorCode
	cause   error
	context string
}

var _ error = (*Error)(nil)

// NewError creates new instance of domain error
func NewError(code ErrorCode, cause error) *Error {
	return &Error{
		Code:  code,
		cause: cause,
	}
}

// ErrorFrom creates new [Error] from provided cause
func ErrorFrom(cause error) *Error {
	return &Error{
		cause: cause,
	}
}

// Error implements [error]
func (e *Error) Error() string {
	if e.cause == nil {
		return e.UserMessage()
	}
	return fmt.Sprintf("%s, caused by %s", e.UserMessage(), e.cause.Error())
}

// Unwrap implements standard error unwrapping mechanism
func (e *Error) Unwrap() error {
	return e.cause
}

// WithContext set additional context to [Error]
func (e *Error) WithContext(context string) *Error {
	e.context = context
	return e
}

// UserMessage returns user-facing error representation
func (e *Error) UserMessage() string {
	if e.context == "" {
		return e.Code.String()
	}
	return fmt.Sprintf("%s: %s", e.Code.String(), e.context)
}
