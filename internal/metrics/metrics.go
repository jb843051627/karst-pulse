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
	r.mu.Lock()
	r.counters[name] += value
	r.mu.Unlock()
}

func (r *Registry) Value(name string) int64 {
	if r == nil {
		return 0
	}
	return r.counters[name]
}

func (r *Registry) Reset(name string) {
	if r == nil {
		return
	}
	r.mu.Lock()
	delete(r.counters, name)
	r.mu.Unlock()
}
