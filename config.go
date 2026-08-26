package main

import (
	"flag"
	"time"
)

type appConfig struct {
	Address      string
	DatabasePath string
	WebDirectory string
	StartedAt    time.Time
}

func parseConfig() appConfig {
	address := flag.String("addr", ":8080", "HTTP listen address")
	database := flag.String("db", "data/karst-pulse.db", "SQLite file path")
	web := flag.String("web", "web", "static web directory")
	flag.Parse()
	return appConfig{Address: *address, DatabasePath: *database, WebDirectory: *web, StartedAt: time.Now().UTC()}
}
