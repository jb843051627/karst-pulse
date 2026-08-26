package model

import "time"

type SensorKind string

const (
	SensorFlow         SensorKind = "flow"
	SensorTemperature  SensorKind = "temperature"
	SensorConductivity SensorKind = "conductivity"
	SensorLevel        SensorKind = "level"
)

type SensorStatus string

const (
	SensorOnline  SensorStatus = "online"
	SensorOffline SensorStatus = "offline"
	SensorService SensorStatus = "service"
)

type Sensor struct {
	ID            ID           `json:"id"`
	SpringID      ID           `json:"spring_id"`
	SerialNo      string       `json:"serial_no"`
	Kind          SensorKind   `json:"kind"`
	Unit          string       `json:"unit"`
	ThresholdLow  float64      `json:"threshold_low"`
	ThresholdHigh float64      `json:"threshold_high"`
	Status        SensorStatus `json:"status"`
	LastValue     *float64     `json:"last_value,omitempty"`
	LastReadingAt *time.Time   `json:"last_reading_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}

type SensorInput struct {
	SpringID      ID      `json:"spring_id"`
	SerialNo      string  `json:"serial_no"`
	Kind          string  `json:"kind"`
	Unit          string  `json:"unit"`
	ThresholdLow  float64 `json:"threshold_low"`
	ThresholdHigh float64 `json:"threshold_high"`
	Status        string  `json:"status"`
}

func (s SensorInput) WithDefaults() SensorInput {
	if s.Status == "" {
		s.Status = string(SensorOnline)
	}
	return s
}

func (s Sensor) IsHealthy() bool {
	return s.Status == SensorOnline
}
