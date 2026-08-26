package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, destination any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is required")
	}
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("request body must contain one JSON value")
	}
	return nil
}

func DecodeOptionalJSON(r *http.Request, destination any) error {
	if r.Body == nil {
		return nil
	}
	return DecodeJSON(r, destination)
}
