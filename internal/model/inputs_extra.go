package model

import (
	"strings"
	"time"
)

type AlertInput struct {
	SpringID    ID        `json:"spring_id"`
	SensorID    *ID       `json:"sensor_id,omitempty"`
	EventID     *ID       `json:"event_id,omitempty"`
	Kind        string    `json:"kind"`
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	TriggeredAt time.Time `json:"triggered_at"`
}

type AnalysisQuery struct {
	SpringID ID
	From     time.Time
	To       time.Time
}

func (s SpringInput) Normalize() SpringInput {
	s.Code = strings.ToUpper(strings.TrimSpace(s.Code))
	s.Name = strings.TrimSpace(s.Name)
	s.Region = strings.TrimSpace(s.Region)
	s.Aquifer = strings.TrimSpace(s.Aquifer)
	s.Status = CanonicalStatus(s.Status)
	return s.WithDefaults()
}

func (s SensorInput) Normalize() SensorInput {
	s.SerialNo = strings.ToUpper(strings.TrimSpace(s.SerialNo))
	s.Kind = CanonicalStatus(s.Kind)
	s.Unit = strings.TrimSpace(s.Unit)
	s.Status = CanonicalStatus(s.Status)
	return s.WithDefaults()
}

func (r ReadingInput) Normalize() ReadingInput {
	r.Quality = CanonicalStatus(r.Quality)
	r.Source = strings.TrimSpace(r.Source)
	return r
}

func (b BatchInput) Normalize() BatchInput {
	b.BatchCode = strings.ToUpper(strings.TrimSpace(b.BatchCode))
	b.Collector = strings.TrimSpace(b.Collector)
	b.Status = CanonicalStatus(b.Status)
	b.Notes = strings.TrimSpace(b.Notes)
	return b.WithDefaults(time.Now().UTC())
}

func (m MaintenanceInput) Normalize() MaintenanceInput {
	m.Title = strings.TrimSpace(m.Title)
	m.Status = CanonicalStatus(m.Status)
	m.Notes = strings.TrimSpace(m.Notes)
	if m.Status == "" {
		m.Status = string(MaintenancePlanned)
	}
	return m
}

func (a AlertInput) Normalize(now time.Time) Alert {
	triggered := a.TriggeredAt
	if triggered.IsZero() {
		triggered = now
	}
	level := AlertLevel(CanonicalStatus(a.Level))
	if !level.Valid() {
		level = AlertInfo
	}
	return Alert{
		SpringID: a.SpringID, SensorID: a.SensorID, EventID: a.EventID,
		Kind: strings.TrimSpace(a.Kind), Level: level, Status: AlertOpen,
		Message: strings.TrimSpace(a.Message), TriggeredAt: triggered,
	}
}

func (r ReadingInput) HasTimestamp() bool {
	return !r.ObservedAt.IsZero()
}

func (r ReadingInput) IsRecent(now time.Time, window time.Duration) bool {
	if r.ObservedAt.IsZero() {
		return false
	}
	age := now.Sub(r.ObservedAt)
	return age >= 0 && age <= window
}

func (b SamplingBatch) CanAddSample() bool {
	return b.Status == BatchOpen || b.Status == BatchSubmitted
}

func (m MaintenanceTask) CanComplete() bool {
	return m.Status == MaintenancePlanned || m.Status == MaintenanceInProgress || m.Status == MaintenanceOverdue
}

func (a Alert) CanAcknowledge() bool {
	return a.Status == AlertOpen
}
