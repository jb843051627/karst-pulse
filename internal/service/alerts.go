package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (a *App) CreateAlert(ctx context.Context, alert model.Alert) (model.Alert, error) {
	if alert.SpringID <= 0 || alert.Message == "" {
		return model.Alert{}, fmt.Errorf("alert spring_id and message are required")
	}
	if alert.Status == "" {
		alert.Status = model.AlertOpen
	}
	if alert.TriggeredAt.IsZero() {
		alert.TriggeredAt = a.Now()
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	created, err := a.store.CreateAlert(dbctx, alert)
	if err != nil {
		return model.Alert{}, fmt.Errorf("create alert: %w", err)
	}
	a.metrics.Inc("alerts_created_total")
	return created, nil
}

func (a *App) GetAlert(ctx context.Context, id model.ID) (model.Alert, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	alert, err := a.store.GetAlert(dbctx, id)
	if err != nil {
		return model.Alert{}, a.dbError("get alert", err)
	}
	return alert, nil
}

func (a *App) ListAlerts(ctx context.Context, filter model.ListFilter) (model.APIList[model.Alert], error) {
	filter = filter.Normalized()
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	items, err := a.store.ListAlerts(dbctx, filter)
	if err != nil {
		return model.APIList[model.Alert]{}, a.dbError("list alerts", err)
	}
	return page(items, filter), nil
}

func (a *App) AcknowledgeAlert(ctx context.Context, id model.ID) (model.Alert, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	alert, err := a.store.AcknowledgeAlert(dbctx, id, a.Now())
	if err != nil {
		return model.Alert{}, a.dbError("acknowledge alert", err)
	}
	return alert, nil
}
