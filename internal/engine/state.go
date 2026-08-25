package engine

import (
	"sync"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

type springState struct {
	previous float64
	baseline float64
	peak     float64
	lastAt   time.Time
	observed int
	active   bool
}

type stateCache struct {
	mu     sync.RWMutex
	values map[model.ID]springState
}

func newStateCache() *stateCache {
	return &stateCache{values: make(map[model.ID]springState)}
}

func (c *stateCache) get(id model.ID) (springState, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.values[id]
	return value, ok
}

func (c *stateCache) put(id model.ID, value springState) {
	c.mu.Lock()
	c.values[id] = value
	c.mu.Unlock()
}

func (c *stateCache) size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.values)
}
