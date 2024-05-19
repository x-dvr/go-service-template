package main

import (
	"log/slog"
	"math"
	"time"
)

type Config struct {
	Env             string
	Host            string
	Port            string
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

const (
	defaultEnv      = EnvProduction
	defaultPort     = "3000"
	defaultLogLevel = slog.LevelInfo
)

func NewConfig(getenv func(string) string, args []string) *Config {
	var (
		env      string = getenv("APP_ENV")
		port     string = getenv("APP_PORT")
		host     string = getenv("APP_HOST")
		logLevel slog.Level
	)

	if env == "" {
		env = defaultEnv
	}

	if port == "" {
		port = defaultPort
	}

	err := logLevel.UnmarshalText([]byte(getenv("APP_LOG_LEVEL")))
	if err != nil {
		logLevel = defaultLogLevel
	}

	return &Config{
		Env:      env,
		Host:     host,
		Port:     port,
		LogLevel: logLevel,
		// TODO: read from env var
		ShutdownTimeout: 30 * time.Second,
	}
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("env", c.Env),
		slog.String("host", c.Host),
		slog.String("port", c.Port),
		slog.Int("log_level", int(c.LogLevel)),
		slog.Int("shutdown_timeout", int(math.Round((c.ShutdownTimeout.Seconds())))),
	)
}
