package main

import (
	"log/slog"
	"net/http"

	"github.com/x-dvr/go-service-template/echo/core"
	"github.com/x-dvr/go-service-template/echo/out/store"
	xhttp "github.com/x-dvr/go-service-template/http"
	"github.com/x-dvr/go-service-template/logs"
)

func main() {
	logs.Setup()
	echoStore := store.NewMemoryStore()
	echoSvc := core.NewService(echoStore)
	mux := http.NewServeMux()

	RegisterEcho(mux, echoSvc)
	if err := http.ListenAndServe(":8080", xhttp.LoggingMiddleware(mux)); err != nil {
		slog.Error("Service failed", slog.String("error", err.Error()))
	}
}
