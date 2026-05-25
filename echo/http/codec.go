package http

import (
	"encoding/json"
	"net/http"
	"net/url"

	common "github.com/0xde86/go-service-template/core"
	"github.com/0xde86/go-service-template/echo"
)

func DecodeEcho(r *http.Request) (echo.In, error) {
	var body struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return echo.In{}, common.NewError(common.ErrUnprocessable, err)
	}
	u, _ := url.Parse("http://" + r.RemoteAddr)
	return echo.In{
		Data:      body.Data,
		From:      u.Hostname(),
		WithNoise: r.Header.Get("x-with-noise") != "",
		UseCached: r.Header.Get("x-use-cached") != "",
	}, nil
}

func EncodeEcho(w http.ResponseWriter, out echo.Out) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{out.Message})
}
