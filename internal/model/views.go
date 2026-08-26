package model

import (
	"strconv"
	"time"
)

type TrendPoint struct {
	At       time.Time `json:"at"`
	Value    float64   `json:"value"`
	Quality  string    `json:"quality"`
	SensorID ID        `json:"sensor_id"`
}

type PulsePhaseCount struct {
	Phase PulsePhase `json:"phase"`
	Count int        `json:"count"`
}

type AlertLevelCount struct {
	Level AlertLevel `json:"level"`
	Count int        `json:"count"`
}

type SensorCoverage struct {
	SensorID       ID         `json:"sensor_id"`
	SerialNo       string     `json:"serial_no"`
	Kind           SensorKind `json:"kind"`
	ReadingCount   int        `json:"reading_count"`
	LastReadingAt  *time.Time `json:"last_reading_at,omitempty"`
	CoverageStatus string     `json:"coverage_status"`
}

type SpringDashboard struct {
	Spring       Spring            `json:"spring"`
	Analysis     AnalysisSummary   `json:"analysis"`
	Sensors      []SensorCoverage  `json:"sensors"`
	RecentEvents []PulseEvent      `json:"recent_events"`
	OpenAlerts   []Alert           `json:"open_alerts"`
	Maintenance  []MaintenanceTask `json:"maintenance"`
	GeneratedAt  time.Time         `json:"generated_at"`
}

type PulseTimeline struct {
	EventID   ID            `json:"event_id"`
	SpringID  ID            `json:"spring_id"`
	Start     time.Time     `json:"start"`
	Peak      float64       `json:"peak"`
	Duration  time.Duration `json:"duration"`
	Severity  EventSeverity `json:"severity"`
	Completed bool          `json:"completed"`
}

func (p PulseTimeline) DurationLabel() string {
	minutes := int(p.Duration.Round(time.Minute) / time.Minute)
	if minutes < 60 {
		return formatCount(minutes, "分钟")
	}
	hours := minutes / 60
	return formatCount(hours, "小时")
}

func formatCount(value int, unit string) string {
	return strconv.Itoa(value) + unit
}

func (d SpringDashboard) AlertCount() int {
	return len(d.OpenAlerts)
}

func (d SpringDashboard) HasAttention() bool {
	return d.Analysis.WaterSignal == "attention" || len(d.OpenAlerts) > 0
}
