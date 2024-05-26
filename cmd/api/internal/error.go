package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/x-dvr/go-service-template/internal/app"
)

type Error struct {
	Cause   error  `json:"-"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func NewError(status int, message string, cause error) *Error {
	return &Error{
		Cause:   cause,
		Message: message,
		Status:  status,
	}
}

func WrapError(err error) *Error {
	var appError *app.Error
	if errors.As(err, &appError) {
		switch appError.Type {
		case app.ErrNotFound:
			return NewError(
				http.StatusNotFound,
				err.Error(),
				err,
			)
		case app.ErrBadRequest:
			return NewError(
				http.StatusBadRequest,
				err.Error(),
				err,
			)
		}
	}

	return NewError(
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
		err,
	)
}

func (err Error) Error() string {
	return fmt.Sprintf("Api error: %s, status: %d, cause: %s", err.Message, err.Status, err.Cause.Error())
}

func (err Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", err.Message),
		slog.Int("status", err.Status),
		slog.Any("cause", err.Cause),
	)
}
