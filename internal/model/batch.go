package model

import "time"

type BatchStatus string

const (
	BatchOpen      BatchStatus = "open"
	BatchSubmitted BatchStatus = "submitted"
	BatchArchived  BatchStatus = "archived"
)

type SamplingBatch struct {
	ID        ID          `json:"id"`
	SpringID  ID          `json:"spring_id"`
	EventID   *ID         `json:"event_id,omitempty"`
	BatchCode string      `json:"batch_code"`
	SampledAt time.Time   `json:"sampled_at"`
	Collector string      `json:"collector"`
	Status    BatchStatus `json:"status"`
	Notes     string      `json:"notes"`
	CreatedAt time.Time   `json:"created_at"`
}

type BatchInput struct {
	SpringID  ID        `json:"spring_id"`
	EventID   *ID       `json:"event_id,omitempty"`
	BatchCode string    `json:"batch_code"`
	SampledAt time.Time `json:"sampled_at"`
	Collector string    `json:"collector"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes"`
}

type Sample struct {
	ID        ID        `json:"id"`
	BatchID   ID        `json:"batch_id"`
	Parameter string    `json:"parameter"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	CreatedAt time.Time `json:"created_at"`
}

type SampleInput struct {
	Parameter string  `json:"parameter"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}

func (b BatchInput) WithDefaults(now time.Time) BatchInput {
	if b.SampledAt.IsZero() {
		b.SampledAt = now
	}
	if b.Status == "" {
		b.Status = string(BatchOpen)
	}
	return b
}
