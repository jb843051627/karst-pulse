package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(ctx context.Context, path string) (*Store, error) {
	if path == "" {
		path = "data/karst-pulse.db"
	}
	if path == ":memory:" {
		return nil, fmt.Errorf("memory database is not supported")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	store := &Store{db: db}
	if err := store.configure(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.Migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) configure(ctx context.Context) error {
	statements := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA busy_timeout = 5000",
		"PRAGMA journal_mode = WAL",
	}
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("configure sqlite with %q: %w", statement, err)
		}
	}
	return nil
}

func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	return nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

func parseTime(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse timestamp %q: %w", value, err)
	}
	return parsed, nil
}

func parseOptionalTime(value sql.NullString) (*time.Time, error) {
	if !value.Valid || value.String == "" {
		return nil, nil
	}
	parsed, err := parseTime(value.String)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func optionalID(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	result := value.Int64
	return &result
}
