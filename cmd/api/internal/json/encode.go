package json

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/x-dvr/go-service-template/cmd/api/internal"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func EncodeError(w http.ResponseWriter, e *internal.Error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
