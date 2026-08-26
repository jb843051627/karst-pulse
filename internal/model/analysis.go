package model

import "time"

type ReadingStats struct {
	Count int     `json:"count"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
}

type AnalysisSummary struct {
	SpringID        ID           `json:"spring_id"`
	From            time.Time    `json:"from"`
	To              time.Time    `json:"to"`
	Readings        ReadingStats `json:"readings"`
	PulseEvents     int          `json:"pulse_events"`
	OpenAlerts      int          `json:"open_alerts"`
	SamplingBatches int          `json:"sampling_batches"`
	WaterSignal     string       `json:"water_signal"`
}

func (s AnalysisSummary) HasData() bool {
	return s.Readings.Count > 0 || s.PulseEvents > 0 || s.SamplingBatches > 0
}
