package logs

import (
	"log/slog"
	"os"
)

// Setup sets up structured logging
func Setup() {
	h := newHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(h)
	slog.SetDefault(logger)
}
