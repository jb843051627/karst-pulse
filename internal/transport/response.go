package transport

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Data any `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func WriteData(w http.ResponseWriter, status int, value any) {
	WriteJSON(w, status, Envelope{Data: value})
}

func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
