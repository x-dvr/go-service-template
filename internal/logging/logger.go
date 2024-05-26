package logging

import (
	"io"
	"log"
	"log/slog"
)

type Logger struct {
	*slog.Logger
	Legacy   *log.Logger
	logLevel *slog.LevelVar
}

func (l *Logger) SetLogLevel(level slog.Level) {
	l.logLevel.Set(level)
}

type LogConfigurator interface {
	GetLogLevel() slog.Level
	IsDevEnv() bool
}

func NewLogger(cfg LogConfigurator, stdout io.Writer) *Logger {
	logger := &Logger{
		logLevel: new(slog.LevelVar),
	}

	logger.logLevel.Set(cfg.GetLogLevel())

	opts := &slog.HandlerOptions{
		Level: logger.logLevel,
	}

	var handler slog.Handler = slog.NewJSONHandler(stdout, opts)
	if cfg.IsDevEnv() {
		handler = NewDevHandler(stdout, opts)
	}

	logger.Logger = slog.New(NewContextHandler(handler))
	logger.Legacy = slog.NewLogLogger(handler, cfg.GetLogLevel())

	slog.SetDefault(logger.Logger)

	return logger
}
