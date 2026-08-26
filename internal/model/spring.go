package model

import "time"

type SpringStatus string

const (
	SpringActive   SpringStatus = "active"
	SpringInactive SpringStatus = "inactive"
	SpringWatch    SpringStatus = "watch"
)

type Spring struct {
	ID        ID           `json:"id"`
	Code      string       `json:"code"`
	Name      string       `json:"name"`
	Region    string       `json:"region"`
	Aquifer   string       `json:"aquifer"`
	Latitude  float64      `json:"latitude"`
	Longitude float64      `json:"longitude"`
	Status    SpringStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}

type SpringInput struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Aquifer   string  `json:"aquifer"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
}

func (s SpringInput) WithDefaults() SpringInput {
	if s.Status == "" {
		s.Status = string(SpringActive)
	}
	return s
}

func (s Spring) IsOperational() bool {
	return s.Status == SpringActive || s.Status == SpringWatch
}
