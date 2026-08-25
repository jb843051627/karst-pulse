package transport

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func PageFromRequest(r *http.Request) (model.PageInfo, error) {
	limit := 50
	offset := 0
	var err error
	if value := r.URL.Query().Get("limit"); value != "" {
		limit, err = strconv.Atoi(value)
		if err != nil {
			return model.PageInfo{}, fmt.Errorf("limit must be an integer")
		}
	}
	if value := r.URL.Query().Get("offset"); value != "" {
		offset, err = strconv.Atoi(value)
		if err != nil {
			return model.PageInfo{}, fmt.Errorf("offset must be an integer")
		}
	}
	if limit <= 0 || limit > 200 {
		return model.PageInfo{}, fmt.Errorf("limit must be between 1 and 200")
	}
	if offset < 0 {
		return model.PageInfo{}, fmt.Errorf("offset must not be negative")
	}
	return model.PageInfo{Limit: limit, Offset: offset}, nil
}

func ApplyPage[T any](items []T, page model.PageInfo) model.APIList[T] {
	if items == nil {
		items = []T{}
	}
	page.Count = len(items)
	return model.APIList[T]{Items: items, Page: page}
}

func ParseOptionalBool(r *http.Request, name string) (bool, bool, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return false, false, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, true, fmt.Errorf("%s must be true or false", name)
	}
	return parsed, true, nil
}

func QueryString(r *http.Request, name, fallback string) string {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback
	}
	return value
}
