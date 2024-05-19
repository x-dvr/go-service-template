package app

import (
	"fmt"
	"log/slog"
)

type Error struct {
	Cause   error  `json:"-"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (err Error) Error() string {
	return fmt.Sprintf("App error: %s, code: %d, cause: %s", err.Message, err.Code, err.Cause.Error())
}

func NewError(cause error, code int, message string) *Error {
	return &Error{
		Cause:   cause,
		Message: message,
		Code:    code,
	}
}

func (e Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", e.Message),
		slog.Int("code", e.Code),
		slog.Any("cause", e.Cause),
	)
}
