package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/albertogviana/port-service/internal/web"

	"github.com/albertogviana/port-service/internal/port"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	errorMessage = "something went wrong, please try again later"
)

func getPortByUnloc(
	service port.UseCase,
	log *log.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unloc := mux.Vars(r)["unloc"]

		p, err := service.GetPortByUnLoc(context.Background(), unloc)
		if err != nil {
			switch err {
			case port.ErrPortNotFound:
				web.RespondWithError(
					w,
					http.StatusNotFound,
					fmt.Sprintf("port with unloc %s was not found", unloc),
				)
			default:
				log.Error(err)

				web.RespondWithError(
					w,
					http.StatusInternalServerError,
					errorMessage,
				)
			}

			return
		}

		web.Respond(w, http.StatusOK, p)
	})
}

// createOrUpdatePort create or update a port.
func createOrUpdatePort(
	service port.UseCase,
	log *log.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p port.Port
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			log.Error(fmt.Errorf("failed to decode json: %w", err))

			web.RespondWithError(
				w,
				http.StatusBadRequest,
				"failed to decode json",
			)

			return
		}

		err := service.SavePort(context.Background(), &p)
		if err != nil {
			log.Error(err)

			web.RespondWithError(
				w,
				http.StatusInternalServerError,
				errorMessage,
			)

			return
		}

		web.Respond(
			w,
			http.StatusCreated,
			nil,
		)
	})
}

// MakePortHandlers make url handler.
func MakePortHandlers(
	router *mux.Router,
	service port.UseCase,
	log *log.Logger,
) {
	router.Handle("/v1/port", createOrUpdatePort(service, log)).
		Methods(http.MethodPost).
		Name("createOrUpdatePort")
	router.Handle("/v1/port/{unloc:[0-9A-Z]+}", getPortByUnloc(service, log)).
		Methods(http.MethodGet).
		Name("getPortByUnloc")
}
