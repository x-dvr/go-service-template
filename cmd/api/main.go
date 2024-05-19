package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
)

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
	args []string,
) error {
	cfg := NewConfig(getenv, args)
	logger := NewLogger(cfg, stdout)
	rootRouter := NewRouter(ctx, cfg, logger)

	return StartServer(
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
