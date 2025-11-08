package base_response

import (
	"encoding/json"
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

func SendError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(Error(msg)); err != nil {
		http.Error(w, msg, http.StatusInternalServerError)
	}
}
