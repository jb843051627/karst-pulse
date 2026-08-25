package service

import (
	"context"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func (a *App) Dashboard(ctx context.Context, springID model.ID, from, to time.Time) (model.SpringDashboard, error) {
	spring, err := a.GetSpring(ctx, springID)
	if err != nil {
		return model.SpringDashboard{}, fmt.Errorf("dashboard spring: %w", err)
	}
	analysis, err := a.Analyze(ctx, springID, from, to)
	if err != nil {
		return model.SpringDashboard{}, fmt.Errorf("dashboard analysis: %w", err)
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	coverage, err := a.store.SensorCoverage(dbctx, springID, from, to)
	if err != nil {
		return model.SpringDashboard{}, fmt.Errorf("dashboard coverage: %w", err)
	}
	events, err := a.store.ListRecentEvents(dbctx, springID, 10)
	if err != nil {
		return model.SpringDashboard{}, fmt.Errorf("dashboard events: %w", err)
	}
	alerts, err := a.store.ListAlerts(dbctx, model.ListFilter{SpringID: springID, Status: string(model.AlertOpen), Limit: 20})
	if err != nil {
		return model.SpringDashboard{}, fmt.Errorf("dashboard alerts: %w", err)
	}
	maintenance, err := a.store.ListMaintenance(dbctx, model.ListFilter{SpringID: springID, Limit: 20})
	if err != nil {
		return model.SpringDashboard{}, fmt.Errorf("dashboard maintenance: %w", err)
	}
	return model.SpringDashboard{Spring: spring, Analysis: analysis, Sensors: coverage, RecentEvents: events, OpenAlerts: alerts, Maintenance: maintenance, GeneratedAt: a.Now()}, nil
}

func (a *App) RecentDashboard(ctx context.Context, springID model.ID) (model.SpringDashboard, error) {
	to := a.Now()
	return a.Dashboard(ctx, springID, to.Add(-24*time.Hour), to)
}

func (a *App) Trend(ctx context.Context, filter model.ListFilter) ([]model.TrendPoint, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	points, err := a.store.ReadingSeries(dbctx, filter)
	if err != nil {
		return nil, a.dbError("load reading trend", err)
	}
	return points, nil
}

func (a *App) PulseTimeline(ctx context.Context, eventID model.ID) (model.PulseTimeline, error) {
	event, err := a.GetEvent(ctx, eventID)
	if err != nil {
		return model.PulseTimeline{}, err
	}
	end := event.UpdatedAt
	if event.EndedAt != nil {
		end = *event.EndedAt
	}
	return model.PulseTimeline{EventID: event.ID, SpringID: event.SpringID, Start: event.StartedAt, Peak: event.PeakValue, Duration: end.Sub(event.StartedAt), Severity: event.Severity, Completed: event.EndedAt != nil}, nil
}
