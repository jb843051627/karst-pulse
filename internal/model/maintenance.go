package model

import "time"

type MaintenanceStatus string

const (
	MaintenancePlanned    MaintenanceStatus = "planned"
	MaintenanceInProgress MaintenanceStatus = "in_progress"
	MaintenanceDone       MaintenanceStatus = "done"
	MaintenanceOverdue    MaintenanceStatus = "overdue"
)

type MaintenanceTask struct {
	ID          ID                `json:"id"`
	SpringID    ID                `json:"spring_id"`
	SensorID    ID                `json:"sensor_id"`
	Title       string            `json:"title"`
	DueAt       time.Time         `json:"due_at"`
	Status      MaintenanceStatus `json:"status"`
	Notes       string            `json:"notes"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

type MaintenanceInput struct {
	SpringID ID        `json:"spring_id"`
	SensorID ID        `json:"sensor_id"`
	Title    string    `json:"title"`
	DueAt    time.Time `json:"due_at"`
	Status   string    `json:"status"`
	Notes    string    `json:"notes"`
}

func (m MaintenanceTask) IsDue(now time.Time) bool {
	return m.Status == MaintenancePlanned && !m.DueAt.After(now)
}
