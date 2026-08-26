package transport

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func ID(r *http.Request, key string) (model.ID, error) {
	value := r.PathValue(key)
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return 0, BadRequest(fmt.Sprintf("%s must be a positive integer", key), nil)
	}
	return model.ID(parsed), nil
}

func OptionalID(r *http.Request, key string) (model.ID, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return 0, BadRequest(fmt.Sprintf("%s must be a positive integer", key), nil)
	}
	return model.ID(parsed), nil
}

func Filter(r *http.Request) (model.ListFilter, error) {
	filter := model.ListFilter{Limit: 50}
	var err error
	if value := r.URL.Query().Get("spring_id"); value != "" {
		filter.SpringID, err = parseQueryID("spring_id", value)
		if err != nil {
			return model.ListFilter{}, err
		}
	}
	if value := r.URL.Query().Get("sensor_id"); value != "" {
		filter.SensorID, err = parseQueryID("sensor_id", value)
		if err != nil {
			return model.ListFilter{}, err
		}
	}
	if value := r.URL.Query().Get("event_id"); value != "" {
		filter.EventID, err = parseQueryID("event_id", value)
		if err != nil {
			return model.ListFilter{}, err
		}
	}
	filter.Status = r.URL.Query().Get("status")
	filter.Limit, err = parseIntQuery("limit", r.URL.Query().Get("limit"), 50)
	if err != nil {
		return model.ListFilter{}, err
	}
	filter.Offset, err = parseIntQuery("offset", r.URL.Query().Get("offset"), 0)
	if err != nil {
		return model.ListFilter{}, err
	}
	filter.From, err = parseTimeQuery("from", r.URL.Query().Get("from"))
	if err != nil {
		return model.ListFilter{}, err
	}
	filter.To, err = parseTimeQuery("to", r.URL.Query().Get("to"))
	if err != nil {
		return model.ListFilter{}, err
	}
	return filter, nil
}

func parseQueryID(name, value string) (model.ID, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return 0, BadRequest(fmt.Sprintf("%s must be a positive integer", name), nil)
	}
	return model.ID(parsed), nil
}

func parseIntQuery(name, value string, fallback int) (int, error) {
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, BadRequest(fmt.Sprintf("%s must be an integer", name), nil)
	}
	return parsed, nil
}

func parseTimeQuery(name, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, BadRequest(fmt.Sprintf("%s must be RFC3339", name), nil)
	}
	return parsed.UTC(), nil
}
