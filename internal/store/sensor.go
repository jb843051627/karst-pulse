package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) CreateSensor(ctx context.Context, input model.SensorInput, now time.Time) (model.Sensor, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO sensors (spring_id, serial_no, kind, unit, threshold_low, threshold_high, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, input.SpringID, strings.TrimSpace(input.SerialNo), input.Kind,
		strings.TrimSpace(input.Unit), input.ThresholdLow, input.ThresholdHigh, input.Status, formatTime(now))
	if err != nil {
		return model.Sensor{}, fmt.Errorf("create sensor: %w", err)
	}
	id, err := insertedID(result)
	if err != nil {
		return model.Sensor{}, err
	}
	return s.GetSensor(ctx, model.ID(id))
}

func (s *Store) GetSensor(ctx context.Context, id model.ID) (model.Sensor, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, serial_no, kind, unit, threshold_low, threshold_high, status,
		       last_value, last_reading_at, created_at
		FROM sensors WHERE id = ?`, id)
	var sensor model.Sensor
	var recentValue sqlNullFloat
	var recentReading, created sql.NullString
	if err := row.Scan(&sensor.ID, &sensor.SpringID, &sensor.SerialNo, &sensor.Kind, &sensor.Unit,
		&sensor.ThresholdLow, &sensor.ThresholdHigh, &sensor.Status, &recentValue, &recentReading, &created); err != nil {
		return model.Sensor{}, fmt.Errorf("get sensor %d: %w", id, err)
	}
	sensor.LastValue = recentValue.pointer()
	var err error
	if recentReading.Valid && recentReading.String != "" {
		parsed, parseErr := scanTime(recentReading.String)
		if parseErr != nil {
			return model.Sensor{}, fmt.Errorf("parse sensor reading timestamp: %w", parseErr)
		}
		sensor.LastReadingAt = &parsed
	}
	sensor.CreatedAt, err = scanTime(created.String)
	if err != nil {
		return model.Sensor{}, fmt.Errorf("parse sensor created timestamp: %w", err)
	}
	return sensor, nil
}

func (s *Store) ListSensors(ctx context.Context, filter model.ListFilter) ([]model.Sensor, error) {
	filter = filter.Normalized()
	query := `
		SELECT id, spring_id, serial_no, kind, unit, threshold_low, threshold_high, status,
		       last_value, last_reading_at, created_at FROM sensors`
	args := make([]any, 0, 3)
	conditions := make([]string, 0, 2)
	if filter.SpringID > 0 {
		conditions = append(conditions, "spring_id = ?")
		args = append(args, filter.SpringID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}
	query += ` ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list sensors: %w", err)
	}
	defer rows.Close()
	items := make([]model.Sensor, 0)
	for rows.Next() {
		var item model.Sensor
		var recentValue sqlNullFloat
		var recentReading, created sql.NullString
		if err := rows.Scan(&item.ID, &item.SpringID, &item.SerialNo, &item.Kind, &item.Unit,
			&item.ThresholdLow, &item.ThresholdHigh, &item.Status, &recentValue, &recentReading, &created); err != nil {
			return nil, fmt.Errorf("scan sensor: %w", err)
		}
		item.LastValue = recentValue.pointer()
		if recentReading.Valid && recentReading.String != "" {
			parsed, parseErr := scanTime(recentReading.String)
			if parseErr != nil {
				return nil, fmt.Errorf("parse sensor reading timestamp: %w", parseErr)
			}
			item.LastReadingAt = &parsed
		}
		item.CreatedAt, err = scanTime(created.String)
		if err != nil {
			return nil, fmt.Errorf("parse sensor created timestamp: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sensors: %w", err)
	}
	return items, nil
}

func (s *Store) SensorBelongsToSpring(ctx context.Context, sensorID, springID model.ID) (bool, error) {
	var exists int
	if err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sensors WHERE id = ? AND spring_id = ?)`, sensorID, springID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check sensor relation: %w", err)
	}
	return exists == 1, nil
}

func (s *Store) UpdateSensorLatest(ctx context.Context, tx *sql.Tx, sensorID model.ID, value float64, observedAt time.Time) error {
	result, err := tx.ExecContext(ctx, `UPDATE sensors SET last_value = ?, last_reading_at = ?, status = ?
		WHERE id = ? AND (last_reading_at IS NULL OR last_reading_at <= ?)`,
		value, formatTime(observedAt), model.SensorOnline, sensorID, formatTime(observedAt))
	if err != nil {
		return fmt.Errorf("update sensor %d latest reading: %w", sensorID, err)
	}
	count, err := rowsAffected(result)
	if err != nil {
		return err
	}
	if count == 0 {
		var exists int
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sensors WHERE id = ?)`, sensorID).Scan(&exists); err != nil {
			return fmt.Errorf("check sensor %d after stale reading: %w", sensorID, err)
		}
		if exists != 1 {
			return fmt.Errorf("sensor %d not found", sensorID)
		}
	}
	return nil
}

type sqlNullFloat struct {
	value float64
	valid bool
}

func (n *sqlNullFloat) Scan(value any) error {
	if value == nil {
		n.valid = false
		return nil
	}
	switch typed := value.(type) {
	case float64:
		n.value, n.valid = typed, true
	case int64:
		n.value, n.valid = float64(typed), true
	default:
		return fmt.Errorf("unsupported nullable float type %T", value)
	}
	return nil
}

func (n sqlNullFloat) pointer() *float64 {
	if !n.valid {
		return nil
	}
	value := n.value
	return &value
}
