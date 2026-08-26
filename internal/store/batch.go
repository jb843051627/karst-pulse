package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) CreateBatch(ctx context.Context, input model.BatchInput, now time.Time) (model.SamplingBatch, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO sampling_batches (spring_id, event_id, batch_code, sampled_at, collector, status, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, input.SpringID, input.EventID, strings.TrimSpace(input.BatchCode),
		formatTime(input.SampledAt), strings.TrimSpace(input.Collector), input.Status, strings.TrimSpace(input.Notes), formatTime(now))
	if err != nil {
		return model.SamplingBatch{}, fmt.Errorf("create sampling batch: %w", err)
	}
	id, err := insertedID(result)
	if err != nil {
		return model.SamplingBatch{}, err
	}
	return s.GetBatch(ctx, model.ID(id))
}

func (s *Store) GetBatch(ctx context.Context, id model.ID) (model.SamplingBatch, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, event_id, batch_code, sampled_at, collector, status, notes, created_at
		FROM sampling_batches WHERE id = ?`, id)
	return scanBatch(row)
}

func scanBatch(row interface{ Scan(...any) error }) (model.SamplingBatch, error) {
	var batch model.SamplingBatch
	var eventID sql.NullInt64
	var sampled, created string
	if err := row.Scan(&batch.ID, &batch.SpringID, &eventID, &batch.BatchCode, &sampled, &batch.Collector,
		&batch.Status, &batch.Notes, &created); err != nil {
		return model.SamplingBatch{}, fmt.Errorf("scan sampling batch: %w", err)
	}
	if eventID.Valid {
		id := model.ID(eventID.Int64)
		batch.EventID = &id
	}
	var err error
	batch.SampledAt, err = scanTime(sampled)
	if err != nil {
		return model.SamplingBatch{}, err
	}
	batch.CreatedAt, err = scanTime(created)
	if err != nil {
		return model.SamplingBatch{}, err
	}
	return batch, nil
}

func (s *Store) ListBatches(ctx context.Context, filter model.ListFilter) ([]model.SamplingBatch, error) {
	filter = filter.Normalized()
	query := `SELECT id, spring_id, event_id, batch_code, sampled_at, collector, status, notes, created_at FROM sampling_batches`
	args := make([]any, 0, 4)
	conditions := make([]string, 0, 3)
	if filter.SpringID > 0 {
		conditions = append(conditions, "spring_id = ?")
		args = append(args, filter.SpringID)
	}
	if filter.EventID > 0 {
		conditions = append(conditions, "event_id = ?")
		args = append(args, filter.EventID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}
	query += ` ORDER BY sampled_at DESC, id DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list sampling batches: %w", err)
	}
	defer rows.Close()
	items := make([]model.SamplingBatch, 0)
	for rows.Next() {
		item, scanErr := scanBatch(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sampling batches: %w", err)
	}
	return items, nil
}

func (s *Store) CreateSample(ctx context.Context, batchID model.ID, input model.SampleInput, now time.Time) (model.Sample, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO samples (batch_id, parameter, value, unit, created_at) VALUES (?, ?, ?, ?, ?)`,
		batchID, strings.TrimSpace(input.Parameter), input.Value, strings.TrimSpace(input.Unit), formatTime(now))
	if err != nil {
		return model.Sample{}, fmt.Errorf("create sample: %w", err)
	}
	id, err := insertedID(result)
	if err != nil {
		return model.Sample{}, err
	}
	return model.Sample{ID: model.ID(id), BatchID: batchID, Parameter: input.Parameter, Value: input.Value, Unit: input.Unit, CreatedAt: now}, nil
}

func (s *Store) ListSamples(ctx context.Context, batchID model.ID) ([]model.Sample, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, batch_id, parameter, value, unit, created_at FROM samples WHERE batch_id = ? ORDER BY id`, batchID)
	if err != nil {
		return nil, fmt.Errorf("list samples: %w", err)
	}
	defer rows.Close()
	items := make([]model.Sample, 0)
	for rows.Next() {
		var item model.Sample
		var created string
		if err := rows.Scan(&item.ID, &item.BatchID, &item.Parameter, &item.Value, &item.Unit, &created); err != nil {
			return nil, fmt.Errorf("scan sample: %w", err)
		}
		item.CreatedAt, err = scanTime(created)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate samples: %w", err)
	}
	return items, nil
}

func (s *Store) CountBatches(ctx context.Context, springID model.ID, from, to time.Time) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sampling_batches WHERE spring_id = ? AND sampled_at >= ? AND sampled_at <= ?`,
		springID, formatTime(from), formatTime(to)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count batches: %w", err)
	}
	return count, nil
}
