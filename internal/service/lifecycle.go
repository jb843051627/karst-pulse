package service

import (
	"context"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

type SpringReadiness struct {
	SpringID       model.ID `json:"spring_id"`
	SensorCount    int      `json:"sensor_count"`
	RecentReadings int      `json:"recent_readings"`
	Ready          bool     `json:"ready"`
	Reason         string   `json:"reason"`
}

func (a *App) CheckSpringReadiness(ctx context.Context, springID model.ID) (SpringReadiness, error) {
	exists, err := a.SpringExists(ctx, springID)
	if err != nil {
		return SpringReadiness{}, err
	}
	if !exists {
		return SpringReadiness{}, fmt.Errorf("spring %d does not exist", springID)
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	sensorCount, err := a.store.CountSensorsForSpring(dbctx, springID)
	if err != nil {
		return SpringReadiness{}, a.dbError("count readiness sensors", err)
	}
	now := a.Now()
	readingCount, err := a.store.CountReadings(dbctx, springID, now.Add(-24*time.Hour), now)
	if err != nil {
		return SpringReadiness{}, a.dbError("count readiness readings", err)
	}
	result := SpringReadiness{SpringID: springID, SensorCount: sensorCount, RecentReadings: readingCount, Ready: sensorCount > 0}
	if sensorCount == 0 {
		result.Reason = "no sensor is registered"
	} else if readingCount == 0 {
		result.Reason = "no reading arrived in the last 24 hours"
	} else {
		result.Ready = true
		result.Reason = "observation pipeline is active"
	}
	return result, nil
}

func (a *App) AcknowledgeSpringAlerts(ctx context.Context, springID model.ID) (int, error) {
	alerts, err := a.ListAlerts(ctx, model.ListFilter{SpringID: springID, Status: string(model.AlertOpen), Limit: 200})
	if err != nil {
		return 0, err
	}
	count := 0
	for _, alert := range alerts.Items {
		if _, ackErr := a.AcknowledgeAlert(ctx, alert.ID); ackErr != nil {
			return count, fmt.Errorf("acknowledge spring alert %d: %w", alert.ID, ackErr)
		}
		count++
	}
	return count, nil
}

func (a *App) SpringDescription(ctx context.Context, springID model.ID) (string, error) {
	spring, err := a.GetSpring(ctx, springID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s（%s，%s，%s）", spring.Name, spring.Code, spring.Region, spring.Status.Label()), nil
}
