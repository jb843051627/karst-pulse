package engine

type Status struct {
	Started     bool `json:"started"`
	QueueLength int  `json:"queue_length"`
	CacheSize   int  `json:"detector_cache_size"`
}

func (r *Runtime) Status() Status {
	r.mu.Lock()
	started := r.started
	r.mu.Unlock()
	return Status{Started: started, QueueLength: r.QueueLength(), CacheSize: r.detector.CacheSize()}
}
