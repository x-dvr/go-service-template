package main

import (
	"net/http"

	"github.com/0xde86/go-service-template/echo"
	ehttp "github.com/0xde86/go-service-template/echo/http"
	xhttp "github.com/0xde86/go-service-template/http"
)

func RegisterEcho(mux *http.ServeMux, svc *echo.Service) {
	mux.Handle("POST /echo", xhttp.Wrap(ehttp.DecodeEcho, svc.Echo, ehttp.EncodeEcho))
}
