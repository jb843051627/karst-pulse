package validate

import (
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func IDs(values ...model.ID) error {
	var errors Errors
	for index, value := range values {
		if value <= 0 {
			errors.Add("id", "id at index "+indexText(index)+" must be positive")
		}
	}
	return errors.Err()
}

func Coordinates(latitude, longitude float64) error {
	var errors Errors
	Finite(&errors, "latitude", latitude)
	Finite(&errors, "longitude", longitude)
	Range(&errors, "latitude", latitude, -90, 90)
	Range(&errors, "longitude", longitude, -180, 180)
	return errors.Err()
}

func Thresholds(low, high float64) error {
	var errors Errors
	Finite(&errors, "threshold_low", low)
	Finite(&errors, "threshold_high", high)
	if low >= high {
		errors.Add("threshold_high", "must be greater than threshold_low")
	}
	return errors.Err()
}

func NotFuture(errors *Errors, field string, value, now time.Time) {
	if !value.IsZero() && value.After(now.Add(2*time.Minute)) {
		errors.Add(field, "must not be in the future")
	}
}

func Recent(errors *Errors, field string, value, now time.Time, maximumAge time.Duration) {
	if value.IsZero() || maximumAge <= 0 {
		return
	}
	if value.Before(now.Add(-maximumAge)) {
		errors.Add(field, "is outside the accepted time window")
	}
}

func EventWindow(errors *Errors, started, ended *time.Time) {
	if started == nil || ended == nil {
		return
	}
	if ended.Before(*started) {
		errors.Add("ended_at", "must be after started_at")
	}
}

func Page(limit, offset int) error {
	var errors Errors
	if limit < 0 || limit > 200 {
		errors.Add("limit", "must be between 0 and 200")
	}
	if offset < 0 {
		errors.Add("offset", "must not be negative")
	}
	return errors.Err()
}

func indexText(index int) string {
	if index < 10 {
		return string(rune('0' + index))
	}
	return "many"
}
