package metrics

import (
    "sync"
    "testing"
)

func TestBug14_MetricsReadsUseSnapshotLocks(t *testing.T) {
    registry := New()
    var wg sync.WaitGroup
    for worker := 0; worker < 20; worker++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for index := 0; index < 200; index++ {
                registry.Add("http_requests_total", 1)
                _ = registry.Value("http_requests_total")
                _ = registry.Counters()
            }
        }()
    }
    wg.Wait()
}
