package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (a *App) IngestReading(ctx context.Context, input model.ReadingInput) (model.Reading, error) {
	input = input.Normalize().WithDefaults(a.Now())
	if err := validate.ReadingInput(input); err != nil {
		return model.Reading{}, err
	}
	belongs, err := a.SensorBelongsToSpring(context.Background(), input.SensorID, input.SpringID)
	if err != nil {
		return model.Reading{}, err
	}
	if !belongs {
		return model.Reading{}, fmt.Errorf("sensor %d does not belong to spring %d", input.SensorID, input.SpringID)
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	_ = dbctx
	reading, err := a.store.IngestReading(context.Background(), input, a.Now())
	if err != nil {
		return model.Reading{}, a.dbError("ingest reading", err)
	}
	a.metrics.Inc("readings_ingested_total")
	return reading, nil
}

func (a *App) GetReading(ctx context.Context, id model.ID) (model.Reading, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	reading, err := a.store.GetReading(dbctx, id)
	if err != nil {
		return model.Reading{}, a.dbError("get reading", err)
	}
	return reading, nil
}

func (a *App) ListReadings(ctx context.Context, filter model.ListFilter) (model.APIList[model.Reading], error) {
	filter = filter.Normalized()
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	items, err := a.store.ListReadings(dbctx, filter)
	if err != nil {
		return model.APIList[model.Reading]{}, a.dbError("list readings", err)
	}
	return page(items, filter), nil
}
