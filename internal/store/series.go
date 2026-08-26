package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) ReadingSeries(ctx context.Context, filter model.ListFilter) ([]model.TrendPoint, error) {
	filter = filter.Normalized()
	query := `SELECT observed_at, value, quality, sensor_id FROM readings`
	args := make([]any, 0, 5)
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
	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}
	query += ` ORDER BY observed_at ASC, id ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("load reading series: %w", err)
	}
	defer rows.Close()
	items := make([]model.TrendPoint, 0)
	for rows.Next() {
		var item model.TrendPoint
		var observed string
		if err := rows.Scan(&observed, &item.Value, &item.Quality, &item.SensorID); err != nil {
			return nil, fmt.Errorf("scan trend point: %w", err)
		}
		item.At, err = scanTime(observed)
		if err != nil {
			return nil, fmt.Errorf("parse trend timestamp: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reading series: %w", err)
	}
	return items, nil
}

func (s *Store) AverageBySensor(ctx context.Context, springID model.ID, from, to time.Time) ([]model.SensorCoverage, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT s.id, s.serial_no, s.kind, COUNT(r.id), MAX(r.observed_at)
		FROM sensors s LEFT JOIN readings r ON r.sensor_id = s.id
		  AND r.observed_at >= ? AND r.observed_at <= ? AND r.quality != ?
		WHERE s.spring_id = ? GROUP BY s.id, s.serial_no, s.kind ORDER BY s.id`,
		formatTime(from), formatTime(to), model.QualityInvalid, springID)
	if err != nil {
		return nil, fmt.Errorf("load sensor averages: %w", err)
	}
	defer rows.Close()
	items := make([]model.SensorCoverage, 0)
	for rows.Next() {
		var item model.SensorCoverage
		var latest sql.NullString
		if err := rows.Scan(&item.SensorID, &item.SerialNo, &item.Kind, &item.ReadingCount, &latest); err != nil {
			return nil, fmt.Errorf("scan sensor average: %w", err)
		}
		if latest.Valid && latest.String != "" {
			parsed, parseErr := scanTime(latest.String)
			if parseErr != nil {
				return nil, parseErr
			}
			item.LastReadingAt = &parsed
		}
		item.CoverageStatus = "available"
		if item.ReadingCount == 0 {
			item.CoverageStatus = "empty"
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sensor averages: %w", err)
	}
	return items, nil
}
