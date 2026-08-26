package metrics

type Snapshot struct {
	Requests         int64 `json:"requests"`
	RequestErrors    int64 `json:"request_errors"`
	ReadingsIngested int64 `json:"readings_ingested"`
	EventsDetected   int64 `json:"events_detected"`
	AlertsCreated    int64 `json:"alerts_created"`
	WorkerFailures   int64 `json:"worker_failures"`
}

func (r *Registry) Snapshot() Snapshot {
	return Snapshot{
		Requests:         r.counters["http_requests_total"],
		RequestErrors:    r.counters["http_request_errors_total"],
		ReadingsIngested: r.counters["readings_ingested_total"],
		EventsDetected:   r.counters["pulse_events_detected_total"],
		AlertsCreated:    r.counters["alerts_created_total"],
		WorkerFailures:   r.counters["worker_failures_total"],
	}
}
