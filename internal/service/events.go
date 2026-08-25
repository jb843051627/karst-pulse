package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (a *App) CreateEvent(ctx context.Context, evaluation model.PulseEvaluation) (model.PulseEvent, error) {
	exists, err := a.SpringExists(ctx, evaluation.SpringID)
	if err != nil {
		return model.PulseEvent{}, err
	}
	if !exists {
		return model.PulseEvent{}, fmt.Errorf("spring %d does not exist", evaluation.SpringID)
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	event, err := a.store.CreateEvent(dbctx, evaluation, a.Now())
	if err != nil {
		return model.PulseEvent{}, a.dbError("create pulse event", err)
	}
	a.metrics.Inc("pulse_events_detected_total")
	return event, nil
}

func (a *App) GetEvent(ctx context.Context, id model.ID) (model.PulseEvent, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	event, err := a.store.GetEvent(dbctx, id)
	if err != nil {
		return model.PulseEvent{}, a.dbError("get pulse event", err)
	}
	return event, nil
}

func (a *App) ListEvents(ctx context.Context, filter model.ListFilter) (model.APIList[model.PulseEvent], error) {
	filter = filter.Normalized()
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	items, err := a.store.ListEvents(dbctx, filter)
	if err != nil {
		return model.APIList[model.PulseEvent]{}, a.dbError("list pulse events", err)
	}
	return page(items, filter), nil
}

func (a *App) LatestOpenEvent(ctx context.Context, springID model.ID) (*model.PulseEvent, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	event, err := a.store.LatestOpenEvent(dbctx, springID)
	if err != nil {
		return nil, a.dbError("find open pulse event", err)
	}
	return event, nil
}

func (a *App) UpdateEvent(ctx context.Context, event model.PulseEvent) error {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	if err := a.store.UpdateEvent(dbctx, event); err != nil {
		return a.dbError("update pulse event", err)
	}
	return nil
}
