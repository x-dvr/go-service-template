package main

import (
	"net/http"

	"github.com/x-dvr/go-service-template/echo/core"
	ehttp "github.com/x-dvr/go-service-template/echo/in/http"
	xhttp "github.com/x-dvr/go-service-template/http"
)

func RegisterEcho(mux *http.ServeMux, svc *core.Service) {
	mux.Handle("POST /echo", xhttp.Wrap(ehttp.DecodeEcho, svc.Echo, ehttp.EncodeEcho))
}
