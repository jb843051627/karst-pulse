package model

import "time"

type ReadingQuality string

const (
	QualityGood    ReadingQuality = "good"
	QualitySuspect ReadingQuality = "suspect"
	QualityInvalid ReadingQuality = "invalid"
)

type Reading struct {
	ID         ID             `json:"id"`
	SpringID   ID             `json:"spring_id"`
	SensorID   ID             `json:"sensor_id"`
	Value      float64        `json:"value"`
	ObservedAt time.Time      `json:"observed_at"`
	Quality    ReadingQuality `json:"quality"`
	Source     string         `json:"source"`
	CreatedAt  time.Time      `json:"created_at"`
}

type ReadingInput struct {
	SpringID   ID        `json:"spring_id"`
	SensorID   ID        `json:"sensor_id"`
	Value      float64   `json:"value"`
	ObservedAt time.Time `json:"observed_at"`
	Quality    string    `json:"quality"`
	Source     string    `json:"source"`
}

func (r ReadingInput) WithDefaults(now time.Time) ReadingInput {
	if r.ObservedAt.IsZero() {
		r.ObservedAt = now
	}
	if r.Quality == "" {
		r.Quality = string(QualityGood)
	}
	if r.Source == "" {
		r.Source = "sensor"
	}
	return r
}

func (r Reading) IsUsable() bool {
	return r.Quality == QualityGood || r.Quality == QualitySuspect
}
