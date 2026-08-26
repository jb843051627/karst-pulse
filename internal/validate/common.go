package validate

import (
	"math"
	"strings"
	"time"
)

func Required(errors *Errors, field, value string) {
	if strings.TrimSpace(value) == "" {
		errors.Add(field, "required")
	}
}

func Length(errors *Errors, field, value string, minimum, maximum int) {
	length := len([]rune(strings.TrimSpace(value)))
	if length < minimum {
		errors.Add(field, "too short")
	}
	if maximum > 0 && length > maximum {
		errors.Add(field, "too long")
	}
}

func Finite(errors *Errors, field string, value float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		errors.Add(field, "must be finite")
	}
}

func Range(errors *Errors, field string, value, minimum, maximum float64) {
	if value < minimum || value > maximum {
		errors.Add(field, "out of range")
	}
}

func TimePresent(errors *Errors, field string, value time.Time) {
	if value.IsZero() {
		errors.Add(field, "required")
	}
}

func TimeWindow(errors *Errors, from, to time.Time) {
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		errors.Add("from", "must be before to")
	}
}
