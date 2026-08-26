package store

import (
	"context"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (s *Store) Analysis(ctx context.Context, springID model.ID, from, to time.Time) (model.AnalysisSummary, error) {
	stats, err := s.ReadingStats(ctx, springID, from, to)
	if err != nil {
		return model.AnalysisSummary{}, fmt.Errorf("analysis readings: %w", err)
	}
	events, err := s.CountEvents(ctx, springID, from, to)
	if err != nil {
		return model.AnalysisSummary{}, fmt.Errorf("analysis events: %w", err)
	}
	alerts, err := s.CountOpenAlerts(ctx, springID)
	if err != nil {
		return model.AnalysisSummary{}, fmt.Errorf("analysis alerts: %w", err)
	}
	batches, err := s.CountBatches(ctx, springID, from, to)
	if err != nil {
		return model.AnalysisSummary{}, fmt.Errorf("analysis batches: %w", err)
	}
	return model.AnalysisSummary{
		SpringID:        springID,
		From:            from,
		To:              to,
		Readings:        stats,
		PulseEvents:     events,
		OpenAlerts:      alerts,
		SamplingBatches: batches,
		WaterSignal:     waterSignal(stats, events, alerts),
	}, nil
}

func waterSignal(stats model.ReadingStats, events, alerts int) string {
	if alerts > 0 || events > 2 {
		return "attention"
	}
	if stats.Count == 0 {
		return "no_data"
	}
	if stats.Avg == 0 {
		return "low_flow"
	}
	return "stable"
}
