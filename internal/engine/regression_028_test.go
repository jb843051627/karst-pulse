package engine

import (
    "context"
    "testing"
    "time"
)

func TestBug28_RuntimeStartReturnsAfterLaunchingWorkers(t *testing.T) {
    runtime := NewRuntime(nil, Config{MaintenanceTick: time.Hour})
    parent, cancel := context.WithCancel(context.Background())
    defer cancel()
    done := make(chan struct{})
    go func() {
        runtime.Start(parent)
        close(done)
    }()
    select {
    case <-done:
    case <-time.After(300 * time.Millisecond):
        t.Fatal("Runtime.Start blocked in scheduler")
    }
    runtime.Stop()
}
