package store

import (
	"context"
	"fmt"
)

func (s *Store) Migrate(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS springs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			region TEXT NOT NULL,
			aquifer TEXT NOT NULL,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sensors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			spring_id INTEGER NOT NULL REFERENCES springs(id) ON DELETE CASCADE,
			serial_no TEXT NOT NULL UNIQUE,
			kind TEXT NOT NULL,
			unit TEXT NOT NULL,
			threshold_low REAL NOT NULL,
			threshold_high REAL NOT NULL,
			status TEXT NOT NULL,
			last_value REAL,
			last_reading_at TEXT,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS readings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			spring_id INTEGER NOT NULL REFERENCES springs(id) ON DELETE CASCADE,
			sensor_id INTEGER NOT NULL REFERENCES sensors(id) ON DELETE CASCADE,
			value REAL NOT NULL,
			observed_at TEXT NOT NULL,
			quality TEXT NOT NULL,
			source TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS pulse_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			spring_id INTEGER NOT NULL REFERENCES springs(id) ON DELETE CASCADE,
			phase TEXT NOT NULL,
			severity TEXT NOT NULL,
			baseline REAL NOT NULL,
			peak_value REAL NOT NULL,
			started_at TEXT NOT NULL,
			peaked_at TEXT,
			ended_at TEXT,
			summary TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sampling_batches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			spring_id INTEGER NOT NULL REFERENCES springs(id) ON DELETE CASCADE,
			event_id INTEGER REFERENCES pulse_events(id) ON DELETE SET NULL,
			batch_code TEXT NOT NULL UNIQUE,
			sampled_at TEXT NOT NULL,
			collector TEXT NOT NULL,
			status TEXT NOT NULL,
			notes TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS samples (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			batch_id INTEGER NOT NULL REFERENCES sampling_batches(id) ON DELETE CASCADE,
			parameter TEXT NOT NULL,
			value REAL NOT NULL,
			unit TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			spring_id INTEGER NOT NULL REFERENCES springs(id) ON DELETE CASCADE,
			sensor_id INTEGER REFERENCES sensors(id) ON DELETE SET NULL,
			event_id INTEGER REFERENCES pulse_events(id) ON DELETE SET NULL,
			kind TEXT NOT NULL,
			level TEXT NOT NULL,
			status TEXT NOT NULL,
			message TEXT NOT NULL,
			triggered_at TEXT NOT NULL,
			acknowledged_at TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS maintenance_tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			spring_id INTEGER NOT NULL REFERENCES springs(id) ON DELETE CASCADE,
			sensor_id INTEGER NOT NULL REFERENCES sensors(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			due_at TEXT NOT NULL,
			status TEXT NOT NULL,
			notes TEXT NOT NULL,
			completed_at TEXT,
			created_at TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sensors_spring ON sensors(spring_id)`,
		`CREATE INDEX IF NOT EXISTS idx_readings_spring_time ON readings(spring_id, observed_at)`,
		`CREATE INDEX IF NOT EXISTS idx_readings_sensor_time ON readings(sensor_id, observed_at)`,
		`CREATE INDEX IF NOT EXISTS idx_events_spring_time ON pulse_events(spring_id, started_at)`,
		`CREATE INDEX IF NOT EXISTS idx_batches_spring_time ON sampling_batches(spring_id, sampled_at)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status, triggered_at)`,
		`CREATE INDEX IF NOT EXISTS idx_maintenance_due ON maintenance_tasks(status, due_at)`,
	}
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("run schema statement: %w", err)
		}
	}
	return nil
}
