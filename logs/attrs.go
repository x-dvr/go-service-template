package logs

import (
	"context"
	"log/slog"
	"sync"
)

type logFields struct {
	mu    sync.Mutex
	attrs []slog.Attr
}

type ctxKey struct{}

func WithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	lf := &logFields{
		attrs: attrs,
	}
	return context.WithValue(ctx, ctxKey{}, lf)
}

// AddAttrs adds one or more [slog.Attr] to [context.Context]
func AddAttrs(ctx context.Context, attrs ...slog.Attr) {
	if lf, ok := ctx.Value(ctxKey{}).(*logFields); ok {
		lf.mu.Lock()
		lf.attrs = append(lf.attrs, attrs...)
		lf.mu.Unlock()
	}
}

func getAttrs(ctx context.Context) []slog.Attr {
	if lf, ok := ctx.Value(ctxKey{}).(*logFields); ok {
		lf.mu.Lock()
		a := lf.attrs
		lf.mu.Unlock()
		return a
	}
	return nil
}
