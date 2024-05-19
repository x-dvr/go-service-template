package logging

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	slog.Handler
}

func NewContextHandler(h slog.Handler) *ContextHandler {
	return &ContextHandler{h}
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID := GetRequestID(ctx); requestID != "" {
		r.AddAttrs(slog.String("req_id", requestID))
	}

	return h.Handler.Handle(ctx, r)
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}

	return ""
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return ctx
	}

	ctx = context.WithValue(ctx, requestIDKey, requestID)
	return ctx
}

type ctxKey string

var (
	requestIDKey ctxKey = "reqID"
)
