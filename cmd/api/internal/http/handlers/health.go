package handlers

import (
	"net/http"

	"github.com/x-dvr/go-service-template/internal/logging"

	app_http "github.com/x-dvr/go-service-template/cmd/api/internal/http"
	"github.com/x-dvr/go-service-template/cmd/api/internal/json"
)

func HandleGetHealth(
	logger *logging.Logger,
) http.Handler {

	return app_http.NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			return json.Encode(w, http.StatusOK, "OK")
		},
	)
}
