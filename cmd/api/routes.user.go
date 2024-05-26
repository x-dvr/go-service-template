package main

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/x-dvr/go-service-template/cmd/api/internal"
	"github.com/x-dvr/go-service-template/internal/user"
)

func handleGetAllUsers(
	logger *Logger,
	// store *UserStore,
) http.Handler {
	return NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			return Encode(w, http.StatusOK, []user.User{})
		},
	)
}

func handleGetUser(
	logger *Logger,
	// store *UserStore,
) http.Handler {
	return NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				return EncodeError(w, internal.NewError(http.StatusBadRequest, "user id must be integer value", nil))
			}
			user := &user.User{
				ID:   id,
				Name: "John",
			}
			logger.LogAttrs(
				r.Context(),
				slog.LevelDebug,
				"User",
				slog.Any("user", user),
			)
			return Encode(w, http.StatusOK, user)
		},
	)
}

func handleCreateUser(
	logger *Logger,
	// store *UserStore,
) http.Handler {
	type request struct {
		Name string `json:"name"`
	}

	return NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			requestBody, err := Decode[request](r)
			if err != nil {
				return err
			}

			usr := user.User{
				ID:   42,
				Name: requestBody.Name,
			}

			return Encode(w, http.StatusOK, usr)
		},
	)
}
