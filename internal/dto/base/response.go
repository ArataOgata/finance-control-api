package base

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

type Response struct {
	StatusCode string      `json:"status_code"`
	Error      string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "error"
)

func Error(msg string) Response {
	return Response{
		StatusCode: StatusError,
		Error:      msg,
	}
}

func SendJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{
		StatusCode: StatusOK,
		Data:       data,
	}); err != nil {
		log.Println("Error encoding response", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func SendError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(Error(msg)); err != nil {
		http.Error(w, msg, http.StatusInternalServerError)
	}
}
