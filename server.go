package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/metrics"
)

func newServer(address string, registry *metrics.Registry, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           metrics.Middleware(registry, handler),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func serveUntilStopped(ctx context.Context, server *http.Server, logger *log.Logger) error {
	result := make(chan error, 1)
	go func() {
		result <- server.ListenAndServe()
	}()
	select {
	case err := <-result:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
		logger.Printf("开始关闭 HTTP 服务: %v", ctx.Err())
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		return nil
	}
}
