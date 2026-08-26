package model

import "time"

type AlertLevel string

const (
	AlertInfo     AlertLevel = "info"
	AlertWarning  AlertLevel = "warning"
	AlertCritical AlertLevel = "critical"
)

type AlertStatus string

const (
	AlertOpen         AlertStatus = "open"
	AlertAcknowledged AlertStatus = "acknowledged"
)

type Alert struct {
	ID             ID          `json:"id"`
	SpringID       ID          `json:"spring_id"`
	SensorID       *ID         `json:"sensor_id,omitempty"`
	EventID        *ID         `json:"event_id,omitempty"`
	Kind           string      `json:"kind"`
	Level          AlertLevel  `json:"level"`
	Status         AlertStatus `json:"status"`
	Message        string      `json:"message"`
	TriggeredAt    time.Time   `json:"triggered_at"`
	AcknowledgedAt *time.Time  `json:"acknowledged_at,omitempty"`
}

func (a Alert) IsOpen() bool {
	return a.Status == AlertOpen
}

func (a Alert) IsCritical() bool {
	return a.Level == AlertCritical
}
