package app

import (
	"errors"
	"log/slog"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type Error struct {
	Type  error
	Cause error
}

func NewError(appError, cause error) *Error {
	return &Error{
		Type:  appError,
		Cause: cause,
	}
}

func (err Error) Error() string {
	return errors.Join(err.Type, err.Cause).Error()
}

func (err Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("type", err.Type),
		slog.Any("cause", err.Cause),
	)
}
