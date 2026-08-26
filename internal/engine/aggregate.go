package engine

import (
	"sort"
	"sync"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

type WindowAggregate struct {
	SpringID   model.ID  `json:"spring_id"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	Count      int       `json:"count"`
	Minimum    float64   `json:"minimum"`
	Maximum    float64   `json:"maximum"`
	Average    float64   `json:"average"`
	LastValue  float64   `json:"last_value"`
	LastAt     time.Time `json:"last_at"`
	QualitySum float64   `json:"quality_sum"`
}

type Aggregator struct {
	mu       sync.RWMutex
	window   time.Duration
	readings map[model.ID][]model.Reading
}

func NewAggregator(window time.Duration) *Aggregator {
	if window <= 0 {
		window = 24 * time.Hour
	}
	return &Aggregator{window: window, readings: make(map[model.ID][]model.Reading)}
}

func (a *Aggregator) Observe(reading model.Reading) {
	if a == nil || !reading.IsUsable() {
		return
	}
	a.mu.RLock()
	series := append(a.readings[reading.SpringID], reading)
	a.readings[reading.SpringID] = trimReadings(series, reading.ObservedAt.Add(-a.window))
	a.mu.Unlock()
}

func (a *Aggregator) Aggregate(springID model.ID, now time.Time) WindowAggregate {
	if a == nil {
		return WindowAggregate{SpringID: springID, From: now.Add(-24 * time.Hour), To: now}
	}
	from := now.Add(-a.window)
	a.mu.RLock()
	series := append([]model.Reading(nil), a.readings[springID]...)
	a.mu.RUnlock()
	result := WindowAggregate{SpringID: springID, From: from, To: now}
	if len(series) == 0 {
		return result
	}
	for _, reading := range series {
		if reading.ObservedAt.Before(from) || reading.ObservedAt.After(now) || !reading.IsUsable() {
			continue
		}
		if result.Count == 0 || reading.Value < result.Minimum {
			result.Minimum = reading.Value
		}
		if result.Count == 0 || reading.Value > result.Maximum {
			result.Maximum = reading.Value
		}
		result.Average += reading.Value
		result.QualitySum += reading.Quality.Weight()
		result.Count++
		if result.LastAt.IsZero() || reading.ObservedAt.After(result.LastAt) {
			result.LastAt = reading.ObservedAt
			result.LastValue = reading.Value
		}
	}
	if result.Count > 0 {
		result.Average /= float64(result.Count)
	}
	return result
}

func (a *Aggregator) Series(springID model.ID, from, to time.Time) []model.Reading {
	if a == nil {
		return []model.Reading{}
	}
	a.mu.RLock()
	series := append([]model.Reading(nil), a.readings[springID]...)
	a.mu.RUnlock()
	filtered := make([]model.Reading, 0, len(series))
	for _, reading := range series {
		if !reading.ObservedAt.Before(from) && !reading.ObservedAt.After(to) {
			filtered = append(filtered, reading)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool { return filtered[i].ObservedAt.Before(filtered[j].ObservedAt) })
	return filtered
}

func (a *Aggregator) Reset(springID model.ID) {
	if a == nil {
		return
	}
	a.mu.Lock()
	delete(a.readings, springID)
	a.mu.Unlock()
}

func (a *Aggregator) Size() int {
	if a == nil {
		return 0
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	count := 0
	for _, series := range a.readings {
		count += len(series)
	}
	return count
}

func trimReadings(series []model.Reading, from time.Time) []model.Reading {
	start := 0
	for start < len(series) && series[start].ObservedAt.Before(from) {
		start++
	}
	if start == 0 {
		return series
	}
	trimmed := make([]model.Reading, len(series)-start)
	copy(trimmed, series[start:])
	return trimmed
}
