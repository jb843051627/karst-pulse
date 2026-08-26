package model

import "time"

type PulsePhase string

const (
	PhaseRising    PulsePhase = "rising"
	PhasePeaked    PulsePhase = "peaked"
	PhaseFading    PulsePhase = "fading"
	PhaseConfirmed PulsePhase = "confirmed"
)

type EventSeverity string

const (
	SeverityInfo     EventSeverity = "info"
	SeverityWarning  EventSeverity = "warning"
	SeverityCritical EventSeverity = "critical"
)

type PulseEvent struct {
	ID        ID            `json:"id"`
	SpringID  ID            `json:"spring_id"`
	Phase     PulsePhase    `json:"phase"`
	Severity  EventSeverity `json:"severity"`
	Baseline  float64       `json:"baseline"`
	PeakValue float64       `json:"peak_value"`
	StartedAt time.Time     `json:"started_at"`
	PeakedAt  *time.Time    `json:"peaked_at,omitempty"`
	EndedAt   *time.Time    `json:"ended_at,omitempty"`
	Summary   string        `json:"summary"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type PulseEvaluation struct {
	SpringID  ID
	Phase     PulsePhase
	Severity  EventSeverity
	Baseline  float64
	PeakValue float64
	At        time.Time
	Summary   string
	Existing  *PulseEvent
	EndAt     time.Time
}

func (e PulseEvent) IsOpen() bool {
	return e.Phase != PhaseConfirmed && e.EndedAt == nil
}

func (e PulseEvent) NeedsAttention() bool {
	return e.Severity == SeverityWarning || e.Severity == SeverityCritical
}
