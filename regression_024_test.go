package main

import (
    "context"
    "io"
    "log"
    "net/http"
    "testing"
    "time"
)

func TestBug24_HTTPServerStopsWhenContextEnds(t *testing.T) {
    server := &http.Server{Addr: "127.0.0.1:0", Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {})}
    ctx, cancel := context.WithCancel(context.Background())
    done := make(chan error, 1)
    go func() {
        done <- serveUntilStopped(ctx, server, log.New(io.Discard, "", 0))
    }()
    time.Sleep(100 * time.Millisecond)
    cancel()
    select {
    case err := <-done:
        if err != nil {
            t.Fatal(err)
        }
    case <-time.After(2 * time.Second):
        t.Fatal("HTTP server did not stop after context cancellation")
    }
}
