package http

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/x-dvr/go-service-template/cmd/api/internal"
	"github.com/x-dvr/go-service-template/cmd/api/internal/config"
	"github.com/x-dvr/go-service-template/cmd/api/internal/json"
	"github.com/x-dvr/go-service-template/internal/logging"
)

func StartServer(
	ctx context.Context,
	cfg *config.Config,
	rootHandler http.Handler,
	logger *logging.Logger,
) error {
	httpServer := &http.Server{
		Addr:     net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:  rootHandler,
		ErrorLog: logger.Legacy,
		BaseContext: func(l net.Listener) context.Context {
			if ctx == nil {
				return context.Background()
			}
			return ctx
		},
	}

	var eg errgroup.Group

	eg.Go(func() error {
		logger.Info(fmt.Sprintf("listening on %s\n", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Error listening and serving: %s\n", err))
			return err
		}
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, cfg.ShutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error(fmt.Sprintf("Error shutting down http server: %s\n", err))
			return err
		}
		return nil
	})

	return eg.Wait()
}

func NewHandler(l *slog.Logger, f handlerFunc) *handler {
	return &handler{
		l: l,
		f: f,
	}
}

type handlerFunc func(http.ResponseWriter, *http.Request) error

type handler struct {
	l *slog.Logger
	f handlerFunc
}

// ServeHTTP calls f(w, r).
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.f(w, r)
	if err == nil {
		return
	}

	h.l.ErrorContext(r.Context(), "Failed to handle request", "error", err)
	json.EncodeError(w, internal.WrapError(err))
}
