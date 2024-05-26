package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/x-dvr/go-service-template/internal/logging"

	"github.com/x-dvr/go-service-template/cmd/api/internal"
	"github.com/x-dvr/go-service-template/cmd/api/internal/config"
	app_http "github.com/x-dvr/go-service-template/cmd/api/internal/http"
	"github.com/x-dvr/go-service-template/cmd/api/internal/http/handlers"
	"github.com/x-dvr/go-service-template/cmd/api/internal/http/middleware"
	"github.com/x-dvr/go-service-template/cmd/api/internal/json"
)

func NewRouter(
	ctx context.Context,
	cfg *config.Config,
	logger *logging.Logger,
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
	logger *logging.Logger,
	cfg *config.Config,
) {
	logger.LogAttrs(ctx, slog.LevelDebug, "Setting up router", slog.Any("config", cfg))

	mux.Handle("GET /health", handlers.HandleGetHealth(logger))
	mux.Handle("GET /api/v1/users/", handlers.HandleGetAllUsers(logger))
	mux.Handle("POST /api/v1/users/", handlers.HandleCreateUser(logger))
	mux.Handle("GET /api/v1/users/{id}", handlers.HandleGetUser(logger))

	mux.Handle("/", handleNotFound(logger))
}

func handleNotFound(logger *logging.Logger) http.Handler {
	return app_http.NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			return json.EncodeError(w, internal.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil))
		},
	)
}
