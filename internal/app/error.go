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

func NewError(cause error, code int, message string) *Error {
	return &Error{
		Cause:   cause,
		Message: message,
		Code:    code,
	}
}

func (err Error) Error() string {
	return fmt.Sprintf("App error: %s, code: %d, cause: %s", err.Message, err.Code, err.Cause.Error())
}

func (err Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", err.Message),
		slog.Int("code", err.Code),
		slog.Any("cause", err.Cause),
	)
}
