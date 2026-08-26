package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) IngestReading(ctx context.Context, input model.ReadingInput, now time.Time) (model.Reading, error) {
	var reading model.Reading
	err := s.transaction(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `
			INSERT INTO readings (spring_id, sensor_id, value, observed_at, quality, source, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`, input.SpringID, input.SensorID, input.Value,
			formatTime(input.ObservedAt), input.Quality, strings.TrimSpace(input.Source), formatTime(now))
		if err != nil {
			return fmt.Errorf("insert reading: %w", err)
		}
		id, err := insertedID(result)
		if err != nil {
			return err
		}
		if err := s.UpdateSensorLatest(ctx, tx, input.SensorID, input.Value, input.ObservedAt); err != nil {
			return err
		}
		reading = model.Reading{
			ID:         model.ID(id),
			SpringID:   input.SpringID,
			SensorID:   input.SensorID,
			Value:      input.Value,
			ObservedAt: input.ObservedAt,
			Quality:    model.ReadingQuality(input.Quality),
			Source:     input.Source,
			CreatedAt:  now,
		}
		return nil
	})
	if err != nil {
		return model.Reading{}, fmt.Errorf("ingest reading: %w", err)
	}
	return reading, nil
}

func (s *Store) GetReading(ctx context.Context, id model.ID) (model.Reading, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, sensor_id, value, observed_at, quality, source, created_at
		FROM readings WHERE id = ?`, id)
	return scanReading(row)
}

func scanReading(row interface{ Scan(...any) error }) (model.Reading, error) {
	var reading model.Reading
	var observed, created string
	if err := row.Scan(&reading.ID, &reading.SpringID, &reading.SensorID, &reading.Value,
		&observed, &reading.Quality, &reading.Source, &created); err != nil {
		return model.Reading{}, fmt.Errorf("scan reading: %w", err)
	}
	var err error
	reading.ObservedAt, err = scanTime(observed)
	if err != nil {
		return model.Reading{}, fmt.Errorf("parse observed timestamp: %w", err)
	}
	reading.CreatedAt, err = scanTime(created)
	if err != nil {
		return model.Reading{}, fmt.Errorf("parse created timestamp: %w", err)
	}
	return reading, nil
}

func (s *Store) ListReadings(ctx context.Context, filter model.ListFilter) ([]model.Reading, error) {
	filter = filter.Normalized()
	query := `SELECT id, spring_id, sensor_id, value, observed_at, quality, source, created_at FROM readings`
	args := make([]any, 0, 6)
	conditions := make([]string, 0, 4)
	if filter.SpringID > 0 {
		conditions = append(conditions, "spring_id = ?")
		args = append(args, filter.SpringID)
	}
	if filter.SensorID > 0 {
		conditions = append(conditions, "sensor_id = ?")
		args = append(args, filter.SensorID)
	}
	if !filter.From.IsZero() {
		conditions = append(conditions, "observed_at >= ?")
		args = append(args, formatTime(filter.From))
	}
	if !filter.To.IsZero() {
		conditions = append(conditions, "observed_at <= ?")
		args = append(args, formatTime(filter.To))
	}
	if filter.QualityCondition() != "" {
		conditions = append(conditions, "quality = ?")
		args = append(args, filter.QualityCondition())
	}
	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}
	query += ` ORDER BY observed_at DESC, id DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list readings: %w", err)
	}
	defer rows.Close()
	items := make([]model.Reading, 0)
	for rows.Next() {
		item, scanErr := scanReading(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate readings: %w", err)
	}
	return items, nil
}

func (s *Store) RecentReadings(ctx context.Context, springID model.ID, since time.Time, limit int) ([]model.Reading, error) {
	if limit <= 0 {
		limit = 30
	}
	if limit > 200 {
		limit = 200
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, spring_id, sensor_id, value, observed_at, quality, source, created_at
		FROM readings WHERE spring_id = ? AND observed_at >= ?
		ORDER BY observed_at DESC, id DESC LIMIT ?`, springID, formatTime(since), limit)
	if err != nil {
		return nil, fmt.Errorf("find recent readings: %w", err)
	}
	defer rows.Close()
	items := make([]model.Reading, 0)
	for rows.Next() {
		item, scanErr := scanReading(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent readings: %w", err)
	}
	return items, nil
}

func (s *Store) ReadingStats(ctx context.Context, springID model.ID, from, to time.Time) (model.ReadingStats, error) {
	var stats model.ReadingStats
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(MIN(value), 0), COALESCE(MAX(value), 0), COALESCE(AVG(value), 0)
		FROM readings WHERE spring_id = ? AND observed_at >= ? AND observed_at <= ? AND quality != ?`,
		springID, formatTime(from), formatTime(to), model.QualityInvalid).Scan(&stats.Count, &stats.Min, &stats.Max, &stats.Avg); err != nil {
		return model.ReadingStats{}, fmt.Errorf("calculate reading stats: %w", err)
	}
	return stats, nil
}

func (s *Store) CountReadings(ctx context.Context, springID model.ID, from, to time.Time) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM readings WHERE spring_id = ? AND observed_at >= ? AND observed_at <= ?`,
		springID, formatTime(from), formatTime(to)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count readings: %w", err)
	}
	return count, nil
}
