package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/x-dvr/go-service-template/internal/app"
	"github.com/x-dvr/go-service-template/internal/middleware"
)

func NewRouter(
	ctx context.Context,
	cfg *Config,
	logger *Logger,
	// tenantsStore *TenantsStore,
	// commentsStore *CommentsStore,
	// conversationService *ConversationService,
	// chatGPTService *ChatGPTService,
	// authProxy *authProxy,
) http.Handler {
	mux := http.NewServeMux()
	registerRoutes(
		ctx,
		mux,
		logger,
		cfg,
	)
	var handler http.Handler = mux

	// setup root middleware
	// handler = someMiddleware(handler)
	// handler = someMiddleware2(handler)
	// handler = someMiddleware3(handler)
	handler = middleware.Use(
		middleware.NewRequestLogger(logger.Logger),
		middleware.Recover(logger.Logger),
	).For(handler)

	return handler
}

func registerRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	logger *Logger,
	cfg *Config,
) {
	logger.LogAttrs(ctx, slog.LevelDebug, "Setting up router", slog.Any("config", cfg))

	mux.Handle("GET /health", handleGetHealth(logger))
	mux.Handle("GET /api/v1/users/", handleGetAllUsers(logger))
	mux.Handle("POST /api/v1/users/", handleCreateUser(logger))
	mux.Handle("GET /api/v1/users/{id}", handleGetUser(logger))

	mux.Handle("/", handleNotFound(logger))
}

func handleNotFound(logger *Logger) http.Handler {
	return NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			return EncodeError(w, app.NewError(nil, http.StatusNotFound, http.StatusText(http.StatusNotFound)))
		},
	)
}
