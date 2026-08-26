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
		Requests:         r.Value("http_requests_total"),
		RequestErrors:    r.Value("http_request_errors_total"),
		ReadingsIngested: r.Value("readings_ingested_total"),
		EventsDetected:   r.Value("pulse_events_detected_total"),
		AlertsCreated:    r.Value("alerts_created_total"),
		WorkerFailures:   r.Value("worker_failures_total"),
	}
}
