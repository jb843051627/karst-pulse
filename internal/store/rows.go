package store

import (
	"database/sql"
	"fmt"
	"time"
)

func scanString(row *sql.Row, destination ...any) error {
	if err := row.Scan(destination...); err != nil {
		return fmt.Errorf("scan row: %w", err)
	}
	return nil
}

func scanTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}
	return parseTime(value)
}

func readNullFloat(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}
	result := value.Float64
	return &result
}

func closeRows(rows *sql.Rows) error {
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close rows: %w", err)
	}
	return nil
}
