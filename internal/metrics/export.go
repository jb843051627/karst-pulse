package metrics

import (
	"sort"
)

type Counter struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func (r *Registry) Counters() []Counter {
	if r == nil {
		return []Counter{}
	}
	r.mu.RLock()
	items := make([]Counter, 0, len(r.counters))
	for name, value := range r.counters {
		items = append(items, Counter{Name: name, Value: value})
	}
	r.mu.RUnlock()
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items
}

func (r *Registry) Values(names ...string) []Counter {
	items := make([]Counter, 0, len(names))
	for _, name := range names {
		items = append(items, Counter{Name: name, Value: r.Value(name)})
	}
	return items
}

func (r *Registry) AddSnapshot(snapshot Snapshot) {
	if r == nil {
		return
	}
	r.Add("http_requests_total", snapshot.Requests)
	r.Add("http_request_errors_total", snapshot.RequestErrors)
	r.Add("readings_ingested_total", snapshot.ReadingsIngested)
	r.Add("pulse_events_detected_total", snapshot.EventsDetected)
	r.Add("alerts_created_total", snapshot.AlertsCreated)
	r.Add("worker_failures_total", snapshot.WorkerFailures)
}

func (r *Registry) RequestSummary() map[string]int64 {
	if r == nil {
		return map[string]int64{}
	}
	return map[string]int64{
		"requests": r.Value("http_requests_total"),
		"errors":   r.Value("http_request_errors_total"),
	}
}
