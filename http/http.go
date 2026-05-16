package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/x-dvr/go-service-template/core"
	"github.com/x-dvr/go-service-template/logs"
)

type validator interface{ Validate() error }

// Wrap turns a domain function into an HTTP handler. The decode and encode
// callbacks hold HTTP-specific mappers.
func Wrap[In, Out any](
	decode func(*http.Request) (In, error),
	fn func(context.Context, In) (Out, error),
	encode func(http.ResponseWriter, Out) error,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := decode(r)
		if err != nil {
			writeErr(w, err)
			return
		}

		if v, ok := any(in).(validator); ok {
			if err := v.Validate(); err != nil {
				writeErr(w, err)
				return
			}
		}

		out, err := fn(r.Context(), in)
		if err != nil {
			logs.AddAttrs(r.Context(), slog.String("error", err.Error()))
			writeErr(w, err)
			return
		}

		if err := encode(w, out); err != nil {
			slog.Error("Failed to encode response", slog.String("error", err.Error()))
		}
	})
}

func writeErr(w http.ResponseWriter, err error) {
	var domainErr *core.Error
	if !errors.As(err, &domainErr) {
		domainErr = core.ErrorFrom(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusFor(domainErr.Code))
	_ = json.NewEncoder(w).Encode(map[string]string{"message": domainErr.UserMessage()})
}

func statusFor(c core.ErrorCode) int {
	switch c {
	case core.ErrValidation:
		return http.StatusBadRequest
	case core.ErrUnauthorized:
		return http.StatusUnauthorized
	case core.ErrPermissionDenied:
		return http.StatusForbidden
	case core.ErrNotFound:
		return http.StatusNotFound
	case core.ErrDuplicate:
		return http.StatusConflict
	case core.ErrUnprocessable:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
