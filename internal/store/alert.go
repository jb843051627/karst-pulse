package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) CreateAlert(ctx context.Context, alert model.Alert) (model.Alert, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO alerts (spring_id, sensor_id, event_id, kind, level, status, message, triggered_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, alert.SpringID, alert.SensorID, alert.EventID, alert.Kind,
		alert.Level, alert.Status, alert.Message, formatTime(alert.TriggeredAt))
	if err != nil {
		return model.Alert{}, fmt.Errorf("create alert: %w", err)
	}
	id, err := lastInsertID(result)
	if err != nil {
		return model.Alert{}, err
	}
	alert.ID = model.ID(id)
	return alert, nil
}

func (s *Store) GetAlert(ctx context.Context, id model.ID) (model.Alert, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, sensor_id, event_id, kind, level, status, message, triggered_at, acknowledged_at
		FROM alerts WHERE id = ?`, id)
	return scanAlert(row)
}

func scanAlert(row interface{ Scan(...any) error }) (model.Alert, error) {
	var alert model.Alert
	var sensorID, eventID sql.NullInt64
	var triggered string
	var acknowledged sql.NullString
	if err := row.Scan(&alert.ID, &alert.SpringID, &sensorID, &eventID, &alert.Kind, &alert.Level,
		&alert.Status, &alert.Message, &triggered, &acknowledged); err != nil {
		return model.Alert{}, fmt.Errorf("scan alert: %w", err)
	}
	if sensorID.Valid {
		id := model.ID(sensorID.Int64)
		alert.SensorID = &id
	}
	if eventID.Valid {
		id := model.ID(eventID.Int64)
		alert.EventID = &id
	}
	var err error
	alert.TriggeredAt, err = scanTime(triggered)
	if err != nil {
		return model.Alert{}, err
	}
	alert.AcknowledgedAt, err = parseOptionalTime(acknowledged)
	if err != nil {
		return model.Alert{}, err
	}
	return alert, nil
}

func (s *Store) ListAlerts(ctx context.Context, filter model.ListFilter) ([]model.Alert, error) {
	filter = filter.Normalized()
	query := `SELECT id, spring_id, sensor_id, event_id, kind, level, status, message, triggered_at, acknowledged_at FROM alerts`
	args := make([]any, 0, 4)
	conditions := make([]string, 0, 3)
	if filter.SpringID > 0 {
		conditions = append(conditions, "spring_id = ?")
		args = append(args, filter.SpringID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if !filter.From.IsZero() {
		conditions = append(conditions, "triggered_at >= ?")
		args = append(args, formatTime(filter.From))
	}
	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}
	query += ` ORDER BY triggered_at DESC, id DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list alerts: %w", err)
	}
	defer rows.Close()
	items := make([]model.Alert, 0)
	for rows.Next() {
		item, scanErr := scanAlert(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate alerts: %w", err)
	}
	return items, nil
}

func (s *Store) AcknowledgeAlert(ctx context.Context, id model.ID, now time.Time) (model.Alert, error) {
	result, err := s.db.ExecContext(ctx, `UPDATE alerts SET status = ?, acknowledged_at = ? WHERE id = ? AND status = ?`,
		model.AlertAcknowledged, formatTime(now), id, model.AlertOpen)
	if err != nil {
		return model.Alert{}, fmt.Errorf("acknowledge alert %d: %w", id, err)
	}
	count, err := rowsAffected(result)
	if err != nil {
		return model.Alert{}, err
	}
	if count == 0 {
		return model.Alert{}, fmt.Errorf("alert %d is not open", id)
	}
	return s.GetAlert(ctx, id)
}

func (s *Store) CountOpenAlerts(ctx context.Context, springID model.ID) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE spring_id = ? AND status = ?`, springID, model.AlertOpen).Scan(&count); err != nil {
		return 0, fmt.Errorf("count open alerts: %w", err)
	}
	return count, nil
}
