package store

import (
	"context"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

type SampleStats struct {
	BatchID model.ID
	Count   int
	Minimum float64
	Maximum float64
	Average float64
}

func (s *Store) SetBatchStatus(ctx context.Context, batchID model.ID, status model.BatchStatus) (model.SamplingBatch, error) {
	if !status.Valid() {
		return model.SamplingBatch{}, fmt.Errorf("unsupported batch status %q", status)
	}
	result, err := s.db.ExecContext(ctx, `UPDATE sampling_batches SET status = ? WHERE id = ?`, status, batchID)
	if err != nil {
		return model.SamplingBatch{}, fmt.Errorf("set batch status %d: %w", batchID, err)
	}
	count, err := rowsAffected(result)
	if err != nil {
		return model.SamplingBatch{}, err
	}
	if count == 0 {
		return model.SamplingBatch{}, fmt.Errorf("batch %d not found", batchID)
	}
	return s.GetBatch(ctx, batchID)
}

func (s *Store) SampleStats(ctx context.Context, batchID model.ID) (SampleStats, error) {
	var stats SampleStats
	stats.BatchID = batchID
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(MIN(value), 0), COALESCE(MAX(value), 0), COALESCE(AVG(value), 0)
		FROM samples WHERE batch_id = ?`, batchID).Scan(&stats.Count, &stats.Minimum, &stats.Maximum, &stats.Average); err != nil {
		return SampleStats{}, fmt.Errorf("calculate sample stats: %w", err)
	}
	return stats, nil
}

func (s *Store) LatestBatch(ctx context.Context, springID model.ID) (model.SamplingBatch, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, spring_id, event_id, batch_code, sampled_at, collector, status, notes, created_at
		FROM sampling_batches WHERE spring_id = ? ORDER BY sampled_at DESC, id DESC LIMIT 1`, springID)
	return scanBatch(row)
}

func (s *Store) BatchesWithSamples(ctx context.Context, springID model.ID, from, to time.Time) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT b.id) FROM sampling_batches b
		JOIN samples s ON s.batch_id = b.id
		WHERE b.spring_id = ? AND b.sampled_at >= ? AND b.sampled_at <= ?`, springID, formatTime(from), formatTime(to)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count batches with samples: %w", err)
	}
	return count, nil
}
