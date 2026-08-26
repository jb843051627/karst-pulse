package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) SensorCoverage(ctx context.Context, springID model.ID, from, to time.Time) ([]model.SensorCoverage, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT s.id, s.serial_no, s.kind, COUNT(r.id), MAX(r.observed_at)
		FROM sensors s
		LEFT JOIN readings r ON r.sensor_id = s.id AND r.observed_at >= ? AND r.observed_at <= ?
		WHERE s.spring_id = ?
		GROUP BY s.id, s.serial_no, s.kind
		ORDER BY s.id`, formatTime(from), formatTime(to), springID)
	if err != nil {
		return nil, fmt.Errorf("load sensor coverage: %w", err)
	}
	defer rows.Close()
	items := make([]model.SensorCoverage, 0)
	for rows.Next() {
		var item model.SensorCoverage
		var last sql.NullString
		if err := rows.Scan(&item.SensorID, &item.SerialNo, &item.Kind, &item.ReadingCount, &last); err != nil {
			return nil, fmt.Errorf("scan sensor coverage: %w", err)
		}
		if last.Valid && last.String != "" {
			parsed, parseErr := scanTime(last.String)
			if parseErr != nil {
				return nil, fmt.Errorf("parse coverage timestamp: %w", parseErr)
			}
			item.LastReadingAt = &parsed
		}
		item.CoverageStatus = coverageLabel(item.ReadingCount, item.LastReadingAt, to)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sensor coverage: %w", err)
	}
	return items, nil
}

func coverageLabel(count int, last *time.Time, end time.Time) string {
	if count == 0 || last == nil {
		return "missing"
	}
	if end.Sub(*last) > 6*time.Hour {
		return "stale"
	}
	return "covered"
}

func (s *Store) EventPhaseCounts(ctx context.Context, springID model.ID, from, to time.Time) ([]model.PulsePhaseCount, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT phase, COUNT(*) FROM pulse_events
		WHERE spring_id = ? AND started_at >= ? AND started_at <= ?
		GROUP BY phase ORDER BY phase`, springID, formatTime(from), formatTime(to))
	if err != nil {
		return nil, fmt.Errorf("count event phases: %w", err)
	}
	defer rows.Close()
	items := make([]model.PulsePhaseCount, 0)
	for rows.Next() {
		var item model.PulsePhaseCount
		if err := rows.Scan(&item.Phase, &item.Count); err != nil {
			return nil, fmt.Errorf("scan event phase count: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate event phase counts: %w", err)
	}
	return items, nil
}

func (s *Store) AlertLevelCounts(ctx context.Context, springID model.ID) ([]model.AlertLevelCount, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT level, COUNT(*) FROM alerts WHERE spring_id = ? AND status = ?
		GROUP BY level ORDER BY level`, springID, model.AlertOpen)
	if err != nil {
		return nil, fmt.Errorf("count alert levels: %w", err)
	}
	defer rows.Close()
	items := make([]model.AlertLevelCount, 0)
	for rows.Next() {
		var item model.AlertLevelCount
		if err := rows.Scan(&item.Level, &item.Count); err != nil {
			return nil, fmt.Errorf("scan alert level count: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate alert level counts: %w", err)
	}
	return items, nil
}
