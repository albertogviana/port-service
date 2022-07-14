package web

import (
	"encoding/json"
	"net/http"
)

func Respond(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode == http.StatusNoContent {
		return nil
	}

	if data == nil {
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

// ErrorResponse describes an error message.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) error {
	return Respond(
		w,
		statusCode,
		ErrorResponse{
			Code:    statusCode,
			Message: message,
		},
	)
}
