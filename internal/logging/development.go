package logging

// Inspired by https://github.com/dusted-go/logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"sync"
)

const (
	reset = "\x1b[0m"

	white       = 37
	darkGray    = 90
	lightRed    = 91
	lightYellow = 93
	lightBlue   = 94
	lightCyan   = 96
	lightGray   = 97
)

const (
	timeFormat = "[15:04:05.000]"
)

type DevHandler struct {
	h slog.Handler
	w io.Writer
	b *bytes.Buffer
	m *sync.Mutex
}

func NewDevHandler(w io.Writer, opts *slog.HandlerOptions) *DevHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}
	return &DevHandler{
		b: b,
		w: w,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		m: &sync.Mutex{},
	}
}

func (h *DevHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *DevHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DevHandler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *DevHandler) WithGroup(name string) slog.Handler {
	return &DevHandler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

func (h *DevHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(lightBlue, level)
	case slog.LevelInfo:
		level = colorize(lightCyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	if len(attrs) == 0 {
		fmt.Fprintln(
			h.w,
			colorize(darkGray, r.Time.Format(timeFormat)),
			level,
			colorize(white, r.Message),
		)
		return nil
	}

	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	fmt.Fprintln(
		h.w,
		colorize(darkGray, r.Time.Format(timeFormat)),
		level,
		colorize(white, r.Message),
		colorize(lightGray, string(bytes)),
	)

	return nil
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\x1b[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func (h *DevHandler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}
	return attrs, nil
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}
