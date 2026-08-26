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

func (s *Store) CreateSpring(ctx context.Context, input model.SpringInput, now time.Time) (model.Spring, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO springs (code, name, region, aquifer, latitude, longitude, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		strings.TrimSpace(input.Code), strings.TrimSpace(input.Name), strings.TrimSpace(input.Region),
		strings.TrimSpace(input.Aquifer), input.Latitude, input.Longitude, input.Status, formatTime(now))
	if err != nil {
		return model.Spring{}, fmt.Errorf("create spring: %w", err)
	}
	id, err := insertedID(result)
	if err != nil {
		return model.Spring{}, err
	}
	return s.GetSpring(ctx, model.ID(id))
}

func (s *Store) GetSpring(ctx context.Context, id model.ID) (model.Spring, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, code, name, region, aquifer, latitude, longitude, status, created_at
		FROM springs WHERE id = ?`, id)
	var spring model.Spring
	var created string
	if err := scanString(row, &spring.ID, &spring.Code, &spring.Name, &spring.Region, &spring.Aquifer,
		&spring.Latitude, &spring.Longitude, &spring.Status, &created); err != nil {
		return model.Spring{}, fmt.Errorf("get spring %d: %w", id, err)
	}
	var err error
	spring.CreatedAt, err = scanTime(created)
	if err != nil {
		return model.Spring{}, fmt.Errorf("get spring %d timestamp: %w", id, err)
	}
	return spring, nil
}

func (s *Store) ListSprings(ctx context.Context, filter model.ListFilter) ([]model.Spring, error) {
	filter = filter.Normalized()
	query := `SELECT id, code, name, region, aquifer, latitude, longitude, status, created_at FROM springs`
	args := make([]any, 0, 2)
	if filter.Status != "" {
		query += ` WHERE status = ?`
		args = append(args, filter.Status)
	}
	query += ` ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list springs: %w", err)
	}
	defer rows.Close()
	items := make([]model.Spring, 0)
	for rows.Next() {
		var item model.Spring
		var created string
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Region, &item.Aquifer,
			&item.Latitude, &item.Longitude, &item.Status, &created); err != nil {
			return nil, fmt.Errorf("scan spring: %w", err)
		}
		item.CreatedAt, err = scanTime(created)
		if err != nil {
			return nil, fmt.Errorf("parse spring timestamp: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate springs: %w", err)
	}
	return items, nil
}

func (s *Store) SpringExists(ctx context.Context, id model.ID) (bool, error) {
	var exists int
	if err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM springs WHERE id = ?)`, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("check spring %d: %w", id, err)
	}
	return exists == 1, nil
}

func (s *Store) CountSprings(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM springs`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count springs: %w", err)
	}
	return count, nil
}

func isNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
