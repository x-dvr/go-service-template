package logs

import (
	"context"
	"io"
	"log/slog"
	"sync"
)

type handler struct {
	h slog.Handler
	m *sync.Mutex
}

func newHandler(w io.Writer, opts *slog.HandlerOptions) *handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &handler{
		h: slog.NewJSONHandler(w, opts),
		m: &sync.Mutex{},
	}
}

// Enabled implements [slog.Handler].
func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

// Handle implements [slog.Handler].
func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(getAttrs(ctx)...)
	return h.h.Handle(ctx, r)
}

// WithAttrs implements [slog.Handler].
func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{h: h.h.WithAttrs(attrs), m: h.m}
}

// WithGroup implements [slog.Handler].
func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{h: h.h.WithGroup(name), m: h.m}
}

var _ slog.Handler = (*handler)(nil)
