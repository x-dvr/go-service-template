package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/x-dvr/go-service-template/internal/logging"
	"github.com/x-dvr/go-service-template/internal/user"

	"github.com/x-dvr/go-service-template/cmd/api/internal"
	app_http "github.com/x-dvr/go-service-template/cmd/api/internal/http"
	"github.com/x-dvr/go-service-template/cmd/api/internal/json"
)

func HandleGetAllUsers(
	logger *logging.Logger,
	// store *UserStore,
) http.Handler {
	return app_http.NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			return json.Encode(w, http.StatusOK, []user.User{})
		},
	)
}

func HandleGetUser(
	logger *logging.Logger,
	// store *UserStore,
) http.Handler {
	return app_http.NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				return json.EncodeError(w, internal.NewError(http.StatusBadRequest, "user id must be integer value", nil))
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
			return json.Encode(w, http.StatusOK, user)
		},
	)
}

func HandleCreateUser(
	logger *logging.Logger,
	// store *UserStore,
) http.Handler {
	type request struct {
		Name string `json:"name"`
	}

	return app_http.NewHandler(
		logger.Logger,
		func(w http.ResponseWriter, r *http.Request) error {
			requestBody, err := json.Decode[request](r)
			if err != nil {
				return err
			}

			usr := user.User{
				ID:   42,
				Name: requestBody.Name,
			}

			return json.Encode(w, http.StatusOK, usr)
		},
	)
}
