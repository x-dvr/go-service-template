package main

import (
	"log/slog"
	"net/http"

	"github.com/0xde86/go-service-template/echo"
	"github.com/0xde86/go-service-template/echo/store"
	xhttp "github.com/0xde86/go-service-template/http"
	"github.com/0xde86/go-service-template/logs"
)

func main() {
	logs.Setup()
	echoStore := store.NewInMemory()
	echoSvc := echo.NewService(echoStore)
	mux := http.NewServeMux()

	RegisterEcho(mux, echoSvc)
	if err := http.ListenAndServe(":8080", xhttp.LoggingMiddleware(mux)); err != nil {
		slog.Error("Service failed", slog.String("error", err.Error()))
	}
}
