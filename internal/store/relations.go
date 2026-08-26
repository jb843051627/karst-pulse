package store

import (
	"context"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) relationExists(ctx context.Context, statement string, args ...any) (bool, error) {
	var value int
	if err := s.db.QueryRowContext(ctx, statement, args...).Scan(&value); err != nil {
		return false, fmt.Errorf("check relation: %w", err)
	}
	return value == 1, nil
}

func (s *Store) EventBelongsToSpring(ctx context.Context, eventID, springID model.ID) (bool, error) {
	return s.relationExists(ctx, `SELECT EXISTS(SELECT 1 FROM pulse_events WHERE id = ? AND spring_id = ?)`, eventID, springID)
}

func (s *Store) BatchBelongsToSpring(ctx context.Context, batchID, springID model.ID) (bool, error) {
	return s.relationExists(ctx, `SELECT EXISTS(SELECT 1 FROM sampling_batches WHERE id = ? AND spring_id = ?)`, batchID, springID)
}

func (s *Store) AlertBelongsToSpring(ctx context.Context, alertID, springID model.ID) (bool, error) {
	return s.relationExists(ctx, `SELECT EXISTS(SELECT 1 FROM alerts WHERE id = ? AND spring_id = ?)`, alertID, springID)
}

func (s *Store) MaintenanceBelongsToSpring(ctx context.Context, taskID, springID model.ID) (bool, error) {
	return s.relationExists(ctx, `SELECT EXISTS(SELECT 1 FROM maintenance_tasks WHERE id = ? AND spring_id = ?)`, taskID, springID)
}

func (s *Store) SensorThreshold(ctx context.Context, sensorID model.ID) (float64, float64, model.SensorStatus, error) {
	var low, high float64
	var status model.SensorStatus
	if err := s.db.QueryRowContext(ctx, `SELECT threshold_low, threshold_high, status FROM sensors WHERE id = ?`, sensorID).Scan(&low, &high, &status); err != nil {
		return 0, 0, "", fmt.Errorf("read sensor threshold %d: %w", sensorID, err)
	}
	return low, high, status, nil
}

func (s *Store) LatestReadingForSensor(ctx context.Context, sensorID model.ID) (model.Reading, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, sensor_id, value, observed_at, quality, source, created_at
		FROM readings WHERE sensor_id = ? ORDER BY observed_at DESC, id DESC LIMIT 1`, sensorID)
	return scanReading(row)
}

func (s *Store) LatestReadingForSpring(ctx context.Context, springID model.ID) (model.Reading, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, sensor_id, value, observed_at, quality, source, created_at
		FROM readings WHERE spring_id = ? ORDER BY observed_at DESC, id DESC LIMIT 1`, springID)
	return scanReading(row)
}

func (s *Store) CountSensorsForSpring(ctx context.Context, springID model.ID) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sensors WHERE spring_id = ?`, springID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count spring sensors: %w", err)
	}
	return count, nil
}

func (s *Store) CountReadingsForSensor(ctx context.Context, sensorID model.ID, from, to time.Time) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM readings WHERE sensor_id = ? AND observed_at >= ? AND observed_at <= ?`,
		sensorID, formatTime(from), formatTime(to)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count sensor readings: %w", err)
	}
	return count, nil
}

type QualityCount struct {
	Quality model.ReadingQuality
	Count   int
}

func (s *Store) ReadingQualityCounts(ctx context.Context, springID model.ID, from, to time.Time) ([]QualityCount, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT quality, COUNT(*) FROM readings
		WHERE spring_id = ? AND observed_at >= ? AND observed_at <= ?
		GROUP BY quality ORDER BY quality`, springID, formatTime(from), formatTime(to))
	if err != nil {
		return nil, fmt.Errorf("count reading qualities: %w", err)
	}
	defer rows.Close()
	items := make([]QualityCount, 0)
	for rows.Next() {
		var item QualityCount
		if err := rows.Scan(&item.Quality, &item.Count); err != nil {
			return nil, fmt.Errorf("scan reading quality count: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reading quality counts: %w", err)
	}
	return items, nil
}
