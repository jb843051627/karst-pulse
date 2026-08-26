package engine

import (
    "sync"
    "testing"
    "time"

    "github.com/karst-pulse/karst-pulse/internal/model"
)

func TestBug13_PulseStateCacheProtectsTransitions(t *testing.T) {
    cache := newStateCache()
    var wg sync.WaitGroup
    for worker := 0; worker < 20; worker++ {
        wg.Add(1)
        go func(worker int) {
            defer wg.Done()
            for index := 0; index < 200; index++ {
                cache.put(model.ID(worker%3), springState{previous: float64(index), observedAt: time.Now()})
                _, _ = cache.get(model.ID(index % 3))
            }
        }(worker)
    }
    wg.Wait()
}
