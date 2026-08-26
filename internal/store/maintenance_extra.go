package store

import (
	"context"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

type MaintenanceCount struct {
	Status model.MaintenanceStatus
	Count  int
}

func (s *Store) ListDueMaintenance(ctx context.Context, now time.Time, limit int) ([]model.MaintenanceTask, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, spring_id, sensor_id, title, due_at, status, notes, completed_at, created_at
		FROM maintenance_tasks WHERE status IN (?, ?) AND due_at <= ?
		ORDER BY due_at, id LIMIT ?`, model.MaintenancePlanned, model.MaintenanceOverdue, formatTime(now), limit)
	if err != nil {
		return nil, fmt.Errorf("list due maintenance: %w", err)
	}
	defer rows.Close()
	items := make([]model.MaintenanceTask, 0)
	for rows.Next() {
		item, scanErr := scanMaintenance(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate due maintenance: %w", err)
	}
	return items, nil
}

func (s *Store) SetMaintenanceStatus(ctx context.Context, taskID model.ID, status model.MaintenanceStatus, now time.Time) (model.MaintenanceTask, error) {
	if !status.Valid() {
		return model.MaintenanceTask{}, fmt.Errorf("unsupported maintenance status %q", status)
	}
	var completed any
	if status == model.MaintenanceDone {
		completed = formatTime(now)
	}
	result, err := s.db.ExecContext(ctx, `UPDATE maintenance_tasks SET status = ?, completed_at = ? WHERE id = ?`, status, completed, taskID)
	if err != nil {
		return model.MaintenanceTask{}, fmt.Errorf("set maintenance status %d: %w", taskID, err)
	}
	count, err := rowsAffected(result)
	if err != nil {
		return model.MaintenanceTask{}, err
	}
	if count == 0 {
		return model.MaintenanceTask{}, fmt.Errorf("maintenance task %d not found", taskID)
	}
	return s.GetMaintenance(ctx, taskID)
}

func (s *Store) MaintenanceCounts(ctx context.Context, springID model.ID) ([]MaintenanceCount, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT status, COUNT(*) FROM maintenance_tasks WHERE spring_id = ?
		GROUP BY status ORDER BY status`, springID)
	if err != nil {
		return nil, fmt.Errorf("count maintenance statuses: %w", err)
	}
	defer rows.Close()
	items := make([]MaintenanceCount, 0)
	for rows.Next() {
		var item MaintenanceCount
		if err := rows.Scan(&item.Status, &item.Count); err != nil {
			return nil, fmt.Errorf("scan maintenance count: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate maintenance counts: %w", err)
	}
	return items, nil
}

func (s *Store) OverdueMaintenanceCount(ctx context.Context, springID model.ID, now time.Time) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM maintenance_tasks WHERE spring_id = ?
		AND status IN (?, ?) AND due_at <= ?`, springID, model.MaintenancePlanned, model.MaintenanceOverdue, formatTime(now)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count overdue maintenance: %w", err)
	}
	return count, nil
}
