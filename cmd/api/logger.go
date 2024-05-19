package main

import (
	"io"
	"log"
	"log/slog"

	"github.com/x-dvr/go-service-template/internal/logging"
)

type Logger struct {
	*slog.Logger
	Legacy *log.Logger
}

func NewLogger(cfg *Config, stdout io.Writer) *Logger {
	logLevel := new(slog.LevelVar)
	logLevel.Set(cfg.LogLevel)

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler = slog.NewJSONHandler(stdout, opts)
	if cfg.Env == EnvDevelopment {
		handler = logging.NewDevHandler(stdout, opts)
	}

	logger := slog.New(logging.NewContextHandler(handler))
	slog.SetDefault(logger)

	return &Logger{
		Logger: logger,
		Legacy: slog.NewLogLogger(handler, cfg.LogLevel),
	}
}
