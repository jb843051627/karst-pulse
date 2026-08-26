package validate

import (
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func Filter(filter model.ListFilter) error {
	var errors Errors
	if filter.Limit < 0 || filter.Limit > 200 {
		errors.Add("limit", "must be between 0 and 200")
	}
	if filter.Offset < 0 {
		errors.Add("offset", "must not be negative")
	}
	TimeWindow(&errors, filter.From, filter.To)
	return errors.Err()
}

func AnalysisWindow(from, to time.Time) error {
	var errors Errors
	TimePresent(&errors, "from", from)
	TimePresent(&errors, "to", to)
	TimeWindow(&errors, from, to)
	if !from.IsZero() && !to.IsZero() && to.Sub(from) > 31*24*time.Hour {
		errors.Add("to", "analysis window cannot exceed 31 days")
	}
	return errors.Err()
}
