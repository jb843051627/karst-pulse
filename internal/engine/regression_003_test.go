package engine

import (
    "context"
    "testing"
)

func TestBug03_SchedulerCancellationSignal(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    if !isCanceled(ctx) {
        t.Fatal("scheduler did not observe the canceled parent context")
    }
}
