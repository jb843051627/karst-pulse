package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (a *App) Health(ctx context.Context, startedAt model.Health) (model.Health, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	if err := a.store.Ping(dbctx); err != nil {
		return model.Health{}, fmt.Errorf("database health: %w", err)
	}
	startedAt.Status = "ok"
	startedAt.Database = "sqlite"
	startedAt.Now = a.Now()
	return startedAt, nil
}

func (a *App) MetricsSnapshot() any {
	return a.metrics.Snapshot()
}
