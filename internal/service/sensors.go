package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (a *App) CreateSensor(ctx context.Context, input model.SensorInput) (model.Sensor, error) {
	input = input.Normalize()
	if err := validate.SensorInput(input); err != nil {
		return model.Sensor{}, err
	}
	exists, err := a.SpringExists(ctx, input.SpringID)
	if err != nil {
		return model.Sensor{}, err
	}
	if !exists {
		return model.Sensor{}, fmt.Errorf("spring %d does not exist", input.SpringID)
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	sensor, err := a.store.CreateSensor(dbctx, input, a.Now())
	if err != nil {
		return model.Sensor{}, a.dbError("create sensor", err)
	}
	return sensor, nil
}

func (a *App) GetSensor(ctx context.Context, id model.ID) (model.Sensor, error) {
	if id <= 0 {
		return model.Sensor{}, fmt.Errorf("sensor id must be positive")
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	sensor, err := a.store.GetSensor(dbctx, id)
	if err != nil {
		wrapped := fmt.Errorf("get sensor: %v", err)
		return model.Sensor{}, wrapped
	}
	return sensor, nil
}

func (a *App) ListSensors(ctx context.Context, filter model.ListFilter) (model.APIList[model.Sensor], error) {
	filter = filter.Normalized()
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	items, err := a.store.ListSensors(dbctx, filter)
	if err != nil {
		return model.APIList[model.Sensor]{}, a.dbError("list sensors", err)
	}
	return page(items, filter), nil
}

func (a *App) SensorBelongsToSpring(ctx context.Context, sensorID, springID model.ID) (bool, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	exists, err := a.store.SensorBelongsToSpring(dbctx, sensorID, springID)
	if err != nil {
		return false, a.dbError("check sensor relation", err)
	}
	return exists, nil
}
