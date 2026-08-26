package engine

import (
    "context"
    "sync"
    "testing"
    "time"
)

func TestBug15_RuntimeStartStopSerializesLifecycle(t *testing.T) {
    runtime := NewRuntime(nil, Config{MaintenanceTick: time.Hour})
    var wg sync.WaitGroup
    for index := 0; index < 20; index++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            runtime.Start(context.Background())
        }()
    }
    wg.Wait()
    runtime.Stop()
}
