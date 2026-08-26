package engine

import "time"

type Config struct {
	QueueSize       int
	MaintenanceTick time.Duration
	PulseLift       float64
	PulseDrop       float64
}

func DefaultConfig() Config {
	return Config{
		QueueSize:       64,
		MaintenanceTick: 30 * time.Second,
		PulseLift:       0.15,
		PulseDrop:       0.20,
	}
}

func (c Config) normalized() Config {
	if c.QueueSize < 1 {
		c.QueueSize = 64
	}
	if c.MaintenanceTick <= 0 {
		c.MaintenanceTick = 30 * time.Second
	}
	if c.PulseLift <= 0 {
		c.PulseLift = 0.15
	}
	if c.PulseDrop <= 0 || c.PulseDrop >= 1 {
		c.PulseDrop = 0.20
	}
	return c
}
