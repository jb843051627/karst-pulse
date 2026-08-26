package validate

import "github.com/karst-pulse/karst-pulse/internal/model"

func SpringStatus(value string) error {
	var errors Errors
	if !model.SpringStatus(value).Valid() {
		errors.Add("status", "unsupported spring status")
	}
	return errors.Err()
}

func SensorKind(value string) error {
	var errors Errors
	if !model.SensorKind(value).Valid() {
		errors.Add("kind", "unsupported sensor kind")
	}
	return errors.Err()
}

func SensorStatus(value string) error {
	var errors Errors
	if !model.SensorStatus(value).Valid() {
		errors.Add("status", "unsupported sensor status")
	}
	return errors.Err()
}

func ReadingQuality(value string) error {
	var errors Errors
	if !model.ReadingQuality(value).Valid() {
		errors.Add("quality", "unsupported reading quality")
	}
	return errors.Err()
}

func PulsePhase(value string) error {
	var errors Errors
	if !model.PulsePhase(value).Valid() {
		errors.Add("phase", "unsupported pulse phase")
	}
	return errors.Err()
}

func EventSeverity(value string) error {
	var errors Errors
	if !model.EventSeverity(value).Valid() {
		errors.Add("severity", "unsupported event severity")
	}
	return errors.Err()
}

func BatchStatus(value string) error {
	var errors Errors
	if !model.BatchStatus(value).Valid() {
		errors.Add("status", "unsupported batch status")
	}
	return errors.Err()
}

func AlertValues(level, status string) error {
	var errors Errors
	if !model.AlertLevel(level).Valid() {
		errors.Add("level", "unsupported alert level")
	}
	if !model.AlertStatus(status).Valid() {
		errors.Add("status", "unsupported alert status")
	}
	return errors.Err()
}

func MaintenanceStatus(value string) error {
	var errors Errors
	if !model.MaintenanceStatus(value).Valid() {
		errors.Add("status", "unsupported maintenance status")
	}
	return errors.Err()
}

func StatusCombination(status string, allowed ...string) error {
	var errors Errors
	if !model.IsOneOf(status, allowed...) {
		errors.Add("status", "value is not allowed for this operation")
	}
	return errors.Err()
}
