package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (a *App) CreateBatch(ctx context.Context, input model.BatchInput) (model.SamplingBatch, error) {
	input = input.Normalize().WithDefaults(a.Now())
	if err := validate.BatchInput(input); err != nil {
		return model.SamplingBatch{}, err
	}
	exists, err := a.SpringExists(ctx, input.SpringID)
	if err != nil {
		return model.SamplingBatch{}, err
	}
	if !exists {
		return model.SamplingBatch{}, fmt.Errorf("spring %d does not exist", input.SpringID)
	}
	if input.EventID != nil {
		if _, err := a.GetEvent(ctx, *input.EventID); err != nil {
			return model.SamplingBatch{}, fmt.Errorf("event relation: %w", err)
		}
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	batch, err := a.store.CreateBatch(dbctx, input, a.Now())
	if err != nil {
		return model.SamplingBatch{}, a.dbError("create sampling batch", err)
	}
	return batch, nil
}

func (a *App) GetBatch(ctx context.Context, id model.ID) (model.SamplingBatch, error) {
	if err := ctx.Err(); err != nil {
		return model.SamplingBatch{}, fmt.Errorf("get batch canceled: %w", err)
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	batch, err := a.store.GetBatch(dbctx, id)
	if err != nil {
		return model.SamplingBatch{}, a.dbError("get sampling batch", err)
	}
	return batch, nil
}

func (a *App) ListBatches(ctx context.Context, filter model.ListFilter) (model.APIList[model.SamplingBatch], error) {
	filter = filter.Normalized()
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	items, err := a.store.ListBatches(dbctx, filter)
	if err != nil {
		return model.APIList[model.SamplingBatch]{}, a.dbError("list sampling batches", err)
	}
	return page(items, filter), nil
}

func (a *App) AddSample(ctx context.Context, batchID model.ID, input model.SampleInput) (model.Sample, error) {
	if err := validate.SampleInput(input); err != nil {
		return model.Sample{}, err
	}
	if _, err := a.GetBatch(ctx, batchID); err != nil {
		return model.Sample{}, fmt.Errorf("batch relation: %w", err)
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	sample, err := a.store.CreateSample(dbctx, batchID, input, a.Now())
	if err != nil {
		return model.Sample{}, a.dbError("create sample", err)
	}
	return sample, nil
}

func (a *App) ListSamples(ctx context.Context, batchID model.ID) ([]model.Sample, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	items, err := a.store.ListSamples(dbctx, batchID)
	if err != nil {
		return nil, a.dbError("list samples", err)
	}
	return items, nil
}
