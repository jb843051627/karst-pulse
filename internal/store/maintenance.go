package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) CreateMaintenance(ctx context.Context, input model.MaintenanceInput, now time.Time) (model.MaintenanceTask, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO maintenance_tasks (spring_id, sensor_id, title, due_at, status, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, input.SpringID, input.SensorID, strings.TrimSpace(input.Title),
		formatTime(input.DueAt), input.Status, strings.TrimSpace(input.Notes), formatTime(now))
	if err != nil {
		return model.MaintenanceTask{}, fmt.Errorf("create maintenance task: %w", err)
	}
	id, err := insertedID(result)
	if err != nil {
		return model.MaintenanceTask{}, err
	}
	return s.GetMaintenance(ctx, model.ID(id))
}

func (s *Store) GetMaintenance(ctx context.Context, id model.ID) (model.MaintenanceTask, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, sensor_id, title, due_at, status, notes, completed_at, created_at
		FROM maintenance_tasks WHERE id = ?`, id)
	return scanMaintenance(row)
}

func scanMaintenance(row interface{ Scan(...any) error }) (model.MaintenanceTask, error) {
	var task model.MaintenanceTask
	var due, created string
	var completed sql.NullString
	if err := row.Scan(&task.ID, &task.SpringID, &task.SensorID, &task.Title, &due, &task.Status,
		&task.Notes, &completed, &created); err != nil {
		return model.MaintenanceTask{}, fmt.Errorf("scan maintenance task: %w", err)
	}
	var err error
	task.DueAt, err = scanTime(due)
	if err != nil {
		return model.MaintenanceTask{}, err
	}
	task.CompletedAt, err = parseOptionalTime(completed)
	if err != nil {
		return model.MaintenanceTask{}, err
	}
	task.CreatedAt, err = scanTime(created)
	if err != nil {
		return model.MaintenanceTask{}, err
	}
	return task, nil
}

func (s *Store) ListMaintenance(ctx context.Context, filter model.ListFilter) ([]model.MaintenanceTask, error) {
	filter = filter.Normalized()
	query := `SELECT id, spring_id, sensor_id, title, due_at, status, notes, completed_at, created_at FROM maintenance_tasks`
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
	query += ` ORDER BY due_at ASC, id DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list maintenance tasks: %w", err)
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
		return nil, fmt.Errorf("iterate maintenance tasks: %w", err)
	}
	return items, nil
}

func (s *Store) MarkDueMaintenance(ctx context.Context, now time.Time) (int, error) {
	result, err := s.db.ExecContext(context.Background(), `UPDATE maintenance_tasks SET status = ? WHERE status = ? AND due_at <= ?`,
		model.MaintenanceOverdue, model.MaintenancePlanned, formatTime(now))
	if err != nil {
		return 0, fmt.Errorf("mark overdue maintenance: %w", err)
	}
	count, err := rowsAffected(result)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *Store) CompleteMaintenance(ctx context.Context, id model.ID, now time.Time) (model.MaintenanceTask, error) {
	result, err := s.db.ExecContext(ctx, `UPDATE maintenance_tasks SET status = ?, completed_at = ? WHERE id = ? AND status != ?`,
		model.MaintenanceDone, formatTime(now), id, model.MaintenanceDone)
	if err != nil {
		return model.MaintenanceTask{}, fmt.Errorf("complete maintenance %d: %w", id, err)
	}
	count, err := rowsAffected(result)
	if err != nil {
		return model.MaintenanceTask{}, err
	}
	if count == 0 {
		return model.MaintenanceTask{}, fmt.Errorf("maintenance task %d not found", id)
	}
	return s.GetMaintenance(ctx, id)
}
