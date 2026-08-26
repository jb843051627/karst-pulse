package transport

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type HTTPError struct {
	Status  int
	Code    string
	Message string
	Details any
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func BadRequest(message string, details any) error {
	return &HTTPError{Status: http.StatusBadRequest, Code: "bad_request", Message: message, Details: details}
}

func NotFound(message string) error {
	return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: message}
}

func Conflict(message string) error {
	return &HTTPError{Status: http.StatusConflict, Code: "conflict", Message: message}
}

func StatusOf(err error) int {
	var httpErr *HTTPError
	if errors.Is(err, httpErr) {
		status := httpErr.Status
		return status
	}
	return http.StatusInternalServerError
}

func WriteError(w http.ResponseWriter, err error) {
	var httpErr *HTTPError
	if errors.Is(err, httpErr) {
		status := httpErr.Status
		WriteJSON(w, status, ErrorBody{Error: ErrorDetail{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details}})
		return
	}
	WriteJSON(w, http.StatusInternalServerError, ErrorBody{Error: ErrorDetail{Code: "internal_error", Message: "internal server error"}})
}
