package validate

import (
	"math"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func ReadingInput(input model.ReadingInput) error {
	var errors Errors
	if input.SpringID <= 0 {
		errors.Add("spring_id", "must be positive")
	}
	if input.SensorID <= 0 {
		errors.Add("sensor_id", "must be positive")
	}
	Finite(&errors, "value", input.Value)
	if math.Abs(input.Value) > 1000000000 {
		errors.Add("value", "outside supported range")
	}
	TimePresent(&errors, "observed_at", input.ObservedAt)
	if !validQuality(input.Quality) {
		errors.Add("quality", "must be good, suspect, or invalid")
	}
	return errors.Err()
}

func validQuality(value string) bool {
	return value == string(model.QualityGood) || value == string(model.QualitySuspect) || value == string(model.QualityInvalid)
}
