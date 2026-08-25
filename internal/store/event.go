package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) CreateEvent(ctx context.Context, evaluation model.PulseEvaluation, now time.Time) (model.PulseEvent, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO pulse_events (spring_id, phase, severity, baseline, peak_value, started_at, peaked_at, ended_at, summary, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, evaluation.SpringID, evaluation.Phase, evaluation.Severity,
		evaluation.Baseline, evaluation.PeakValue, formatTime(evaluation.At), nil, nil, evaluation.Summary,
		formatTime(now), formatTime(now))
	if err != nil {
		return model.PulseEvent{}, fmt.Errorf("create pulse event: %w", err)
	}
	id, err := lastInsertID(result)
	if err != nil {
		return model.PulseEvent{}, err
	}
	return s.GetEvent(ctx, model.ID(id))
}

func (s *Store) GetEvent(ctx context.Context, id model.ID) (model.PulseEvent, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, phase, severity, baseline, peak_value, started_at, peaked_at, ended_at, summary, created_at, updated_at
		FROM pulse_events WHERE id = ?`, id)
	return scanEvent(row)
}

func scanEvent(row interface{ Scan(...any) error }) (model.PulseEvent, error) {
	var event model.PulseEvent
	var started, created, updated string
	var peaked, ended sql.NullString
	if err := row.Scan(&event.ID, &event.SpringID, &event.Phase, &event.Severity, &event.Baseline,
		&event.PeakValue, &started, &peaked, &ended, &event.Summary, &created, &updated); err != nil {
		return model.PulseEvent{}, fmt.Errorf("scan pulse event: %w", err)
	}
	var err error
	event.StartedAt, err = scanTime(started)
	if err != nil {
		return model.PulseEvent{}, err
	}
	event.PeakedAt, err = parseOptionalTime(peaked)
	if err != nil {
		return model.PulseEvent{}, err
	}
	event.EndedAt, err = parseOptionalTime(ended)
	if err != nil {
		return model.PulseEvent{}, err
	}
	event.CreatedAt, err = scanTime(created)
	if err != nil {
		return model.PulseEvent{}, err
	}
	event.UpdatedAt, err = scanTime(updated)
	if err != nil {
		return model.PulseEvent{}, err
	}
	return event, nil
}

func (s *Store) ListEvents(ctx context.Context, filter model.ListFilter) ([]model.PulseEvent, error) {
	filter = filter.Normalized()
	query := `SELECT id, spring_id, phase, severity, baseline, peak_value, started_at, peaked_at, ended_at, summary, created_at, updated_at FROM pulse_events`
	args := make([]any, 0, 4)
	conditions := make([]string, 0, 3)
	if filter.SpringID > 0 {
		conditions = append(conditions, "spring_id = ?")
		args = append(args, filter.SpringID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "phase = ?")
		args = append(args, filter.Status)
	}
	if !filter.From.IsZero() {
		conditions = append(conditions, "started_at >= ?")
		args = append(args, formatTime(filter.From))
	}
	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}
	query += ` ORDER BY started_at DESC, id DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list pulse events: %w", err)
	}
	defer rows.Close()
	items := make([]model.PulseEvent, 0)
	for rows.Next() {
		item, scanErr := scanEvent(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pulse events: %w", err)
	}
	return items, nil
}

func (s *Store) LatestOpenEvent(ctx context.Context, springID model.ID) (*model.PulseEvent, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, phase, severity, baseline, peak_value, started_at, peaked_at, ended_at, summary, created_at, updated_at
		FROM pulse_events WHERE spring_id = ? AND ended_at IS NULL AND phase != ?
		ORDER BY started_at DESC, id DESC LIMIT 1`, springID, model.PhaseConfirmed)
	event, err := scanEvent(row)
	if err != nil {
		if errors.Is(unwrapStoreError(err), sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find open event: %w", err)
	}
	return &event, nil
}

func (s *Store) UpdateEvent(ctx context.Context, event model.PulseEvent) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE pulse_events SET phase = ?, severity = ?, baseline = ?, peak_value = ?, started_at = ?, peaked_at = ?, ended_at = ?, summary = ?, updated_at = ? WHERE id = ?`,
		event.Phase, event.Severity, event.Baseline, event.PeakValue, formatTime(event.StartedAt),
		formatOptionalTime(event.PeakedAt), formatOptionalTime(event.EndedAt), event.Summary, formatTime(event.UpdatedAt), event.ID)
	if err != nil {
		return fmt.Errorf("update event %d: %w", event.ID, err)
	}
	count, err := rowsAffected(result)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("event %d not found", event.ID)
	}
	return nil
}

func formatOptionalTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return formatTime(*value)
}

func unwrapStoreError(err error) error {
	return err
}

func (s *Store) CountEvents(ctx context.Context, springID model.ID, from, to time.Time) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pulse_events WHERE spring_id = ? AND started_at >= ? AND started_at <= ?`,
		springID, formatTime(from), formatTime(to)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count events: %w", err)
	}
	return count, nil
}
