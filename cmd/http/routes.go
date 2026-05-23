package main

import (
	"net/http"

	"github.com/x-dvr/go-service-template/echo"
	ehttp "github.com/x-dvr/go-service-template/echo/http"
	xhttp "github.com/x-dvr/go-service-template/http"
)

func RegisterEcho(mux *http.ServeMux, svc *echo.Service) {
	mux.Handle("POST /echo", xhttp.Wrap(ehttp.DecodeEcho, svc.Echo, ehttp.EncodeEcho))
}
