package metrics

import "sync"

type Registry struct {
	mu       sync.RWMutex
	counters map[string]int64
}

func New() *Registry {
	return &Registry{counters: make(map[string]int64)}
}

func (r *Registry) Inc(name string) {
	r.Add(name, 1)
}

func (r *Registry) Add(name string, value int64) {
	if r == nil || name == "" {
		return
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.counters[name] += value
}

func (r *Registry) Value(name string) int64 {
	if r == nil {
		return 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.counters[name]
}

func (r *Registry) Reset(name string) {
	if r == nil {
		return
	}
	r.mu.RLock()
	delete(r.counters, name)
	r.mu.RUnlock()
}
