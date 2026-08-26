package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) FindSpringByCode(ctx context.Context, code string) (model.Spring, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, code, name, region, aquifer, latitude, longitude, status, created_at
		FROM springs WHERE code = ?`, strings.ToUpper(strings.TrimSpace(code)))
	var spring model.Spring
	var created string
	if err := row.Scan(&spring.ID, &spring.Code, &spring.Name, &spring.Region, &spring.Aquifer,
		&spring.Latitude, &spring.Longitude, &spring.Status, &created); err != nil {
		return model.Spring{}, fmt.Errorf("find spring by code: %w", err)
	}
	var err error
	spring.CreatedAt, err = scanTime(created)
	if err != nil {
		return model.Spring{}, fmt.Errorf("parse spring code timestamp: %w", err)
	}
	return spring, nil
}

func (s *Store) ListSpringsByRegion(ctx context.Context, region string, limit int) ([]model.Spring, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, code, name, region, aquifer, latitude, longitude, status, created_at
		FROM springs WHERE region = ? ORDER BY name, id LIMIT ?`, strings.TrimSpace(region), limit)
	if err != nil {
		return nil, fmt.Errorf("list springs by region: %w", err)
	}
	defer rows.Close()
	items := make([]model.Spring, 0)
	for rows.Next() {
		var item model.Spring
		var created string
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Region, &item.Aquifer,
			&item.Latitude, &item.Longitude, &item.Status, &created); err != nil {
			return nil, fmt.Errorf("scan regional spring: %w", err)
		}
		item.CreatedAt, err = scanTime(created)
		if err != nil {
			return nil, fmt.Errorf("parse regional spring timestamp: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate regional springs: %w", err)
	}
	return items, nil
}

func (s *Store) CountSpringsByStatus(ctx context.Context, status model.SpringStatus) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM springs WHERE status = ?`, status).Scan(&count); err != nil {
		return 0, fmt.Errorf("count springs by status: %w", err)
	}
	return count, nil
}
