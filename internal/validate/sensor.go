package validate

import (
	"strings"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func SensorInput(input model.SensorInput) error {
	var errors Errors
	if input.SpringID <= 0 {
		errors.Add("spring_id", "must be positive")
	}
	Required(&errors, "serial_no", input.SerialNo)
	Length(&errors, "serial_no", input.SerialNo, 3, 64)
	if !validSensorKind(input.Kind) {
		errors.Add("kind", "unsupported sensor kind")
	}
	Required(&errors, "unit", input.Unit)
	if input.ThresholdLow >= input.ThresholdHigh {
		errors.Add("threshold_high", "must be greater than threshold_low")
	}
	if !validSensorStatus(input.Status) {
		errors.Add("status", "must be online, offline, or service")
	}
	Finite(&errors, "threshold_low", input.ThresholdLow)
	Finite(&errors, "threshold_high", input.ThresholdHigh)
	return errors.Err()
}

func validSensorKind(value string) bool {
	switch strings.ToLower(value) {
	case string(model.SensorFlow), string(model.SensorTemperature), string(model.SensorConductivity), string(model.SensorLevel):
		return true
	default:
		return false
	}
}

func validSensorStatus(value string) bool {
	return value == string(model.SensorOnline) || value == string(model.SensorOffline) || value == string(model.SensorService)
}
