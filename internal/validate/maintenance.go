package validate

import "github.com/karst-pulse/karst-pulse/internal/model"

func MaintenanceInput(input model.MaintenanceInput) error {
	var errors Errors
	if input.SpringID <= 0 {
		errors.Add("spring_id", "must be positive")
	}
	if input.SensorID <= 0 {
		errors.Add("sensor_id", "must be positive")
	}
	Required(&errors, "title", input.Title)
	Length(&errors, "title", input.Title, 3, 160)
	TimePresent(&errors, "due_at", input.DueAt)
	if !validMaintenanceStatus(input.Status) {
		errors.Add("status", "must be planned, in_progress, done, or overdue")
	}
	return errors.Err()
}

func validMaintenanceStatus(value string) bool {
	return value == string(model.MaintenancePlanned) || value == string(model.MaintenanceInProgress) ||
		value == string(model.MaintenanceDone) || value == string(model.MaintenanceOverdue)
}
