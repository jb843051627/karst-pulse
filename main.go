package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/karst-pulse/karst-pulse/internal/clock"
	"github.com/karst-pulse/karst-pulse/internal/engine"
	"github.com/karst-pulse/karst-pulse/internal/handler"
	"github.com/karst-pulse/karst-pulse/internal/metrics"
	"github.com/karst-pulse/karst-pulse/internal/service"
	"github.com/karst-pulse/karst-pulse/internal/store"
)

func main() {
	config := parseConfig()
	logger := log.New(os.Stdout, "karst-pulse ", log.LstdFlags|log.Lmicroseconds)
	if err := run(config, logger); err != nil {
		logger.Printf("服务退出: %v", err)
		os.Exit(1)
	}
}

func run(config appConfig, logger *log.Logger) error {
	rootContext, stop := signalContext()
	defer stop()
	database, err := store.Open(rootContext, config.DatabasePath)
	if err != nil {
		return err
	}
	defer database.Close()
	registry := metrics.New()
	application := service.New(database, clock.Real{}, registry)
	runtime := engine.NewRuntime(application, engine.DefaultConfig())
	runtime.Start(rootContext)
	defer runtime.Stop()
	static := http.FileServer(http.Dir(config.WebDirectory))
	webHandler := handler.New(application, runtime, config.StartedAt, static)
	server := newServer(config.Address, registry, webHandler.Routes())
	logger.Printf("HTTP 服务监听 %s，数据库 %s", config.Address, config.DatabasePath)
	if err := serveUntilStopped(rootContext, server, logger); err != nil {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	return nil
}
