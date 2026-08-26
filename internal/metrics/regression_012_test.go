package metrics

import (
    "sync"
    "testing"
)

func TestBug12_MetricsRegistryProtectsCounters(t *testing.T) {
    registry := New()
    var wg sync.WaitGroup
    for worker := 0; worker < 20; worker++ {
        wg.Add(1)
        go func(worker int) {
            defer wg.Done()
            for index := 0; index < 200; index++ {
                registry.Add("readings", int64(worker+index))
                registry.Reset("errors")
            }
        }(worker)
    }
    wg.Wait()
}
