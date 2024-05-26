package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/x-dvr/go-service-template/internal/logging"

	"github.com/x-dvr/go-service-template/cmd/api/internal/config"
	"github.com/x-dvr/go-service-template/cmd/api/internal/http"
)

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
	args []string,
) error {
	cfg := config.New(getenv, args)
	logger := logging.NewLogger(cfg, stdout)
	rootRouter := NewRouter(ctx, cfg, logger)

	return http.StartServer(
		ctx,
		cfg,
		rootRouter,
		logger,
	)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Getenv, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		cancel()
		os.Exit(1)
	}
}
