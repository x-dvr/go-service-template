package http

import (
	"encoding/json"
	"net/http"
	"net/url"

	common "github.com/x-dvr/go-service-template/core"
	"github.com/x-dvr/go-service-template/echo/core"
)

func DecodeEcho(r *http.Request) (core.EchoIn, error) {
	var body struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return core.EchoIn{}, common.NewError(common.ErrUnprocessable, err)
	}
	u, _ := url.Parse("http://" + r.RemoteAddr)
	return core.EchoIn{
		Data:      body.Data,
		From:      u.Hostname(),
		WithNoise: r.Header.Get("x-with-noise") != "",
		UseCached: r.Header.Get("x-use-cached") != "",
	}, nil
}

func EncodeEcho(w http.ResponseWriter, out core.EchoOut) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{out.Message})
}
