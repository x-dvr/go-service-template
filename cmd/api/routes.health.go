package main

import "net/http"

func handleGetHealth(
	logger *Logger,
) http.Handler {
	return NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			return Encode(w, http.StatusOK, "OK")
		},
	)
}
