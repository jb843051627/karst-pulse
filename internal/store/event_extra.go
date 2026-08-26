package store

import (
	"context"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) ListRecentEvents(ctx context.Context, springID model.ID, limit int) ([]model.PulseEvent, error) {
	filter := model.ListFilter{SpringID: springID, Limit: limit}
	items, err := s.ListEvents(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list recent events: %w", err)
	}
	return items, nil
}

func (s *Store) ListOpenEvents(ctx context.Context, springID model.ID, limit int) ([]model.PulseEvent, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, spring_id, phase, severity, baseline, peak_value, started_at, peaked_at, ended_at, summary, created_at, updated_at
		FROM pulse_events WHERE spring_id = ? AND ended_at IS NULL AND phase != ?
		ORDER BY started_at DESC, id DESC LIMIT ?`, springID, model.PhaseConfirmed, limit)
	if err != nil {
		return nil, fmt.Errorf("list open events: %w", err)
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
		return nil, fmt.Errorf("iterate open events: %w", err)
	}
	return items, nil
}

func (s *Store) CloseEvent(ctx context.Context, eventID model.ID, at time.Time, summary string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE pulse_events SET phase = ?, ended_at = ?, summary = ?, updated_at = ?
		WHERE id = ? AND ended_at IS NULL`, model.PhaseConfirmed, formatTime(at), summary, formatTime(at), eventID)
	if err != nil {
		return fmt.Errorf("close event %d: %w", eventID, err)
	}
	count, err := rowsAffected(result)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("event %d is already closed or missing", eventID)
	}
	return nil
}

func (s *Store) EventDuration(ctx context.Context, eventID model.ID) (time.Duration, error) {
	event, err := s.GetEvent(ctx, eventID)
	if err != nil {
		return 0, fmt.Errorf("load event duration: %w", err)
	}
	end := event.UpdatedAt
	if event.EndedAt != nil {
		end = *event.EndedAt
	}
	if end.Before(event.StartedAt) {
		return 0, fmt.Errorf("event %d has invalid time range", eventID)
	}
	return end.Sub(event.StartedAt), nil
}

func (s *Store) CountCriticalEvents(ctx context.Context, springID model.ID, from, to time.Time) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM pulse_events WHERE spring_id = ? AND severity = ?
		AND started_at >= ? AND started_at <= ?`, springID, model.SeverityCritical, formatTime(from), formatTime(to)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count critical events: %w", err)
	}
	return count, nil
}
