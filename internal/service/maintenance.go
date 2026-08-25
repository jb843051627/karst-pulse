package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (a *App) CreateMaintenance(ctx context.Context, input model.MaintenanceInput) (model.MaintenanceTask, error) {
	input = input.Normalize()
	if input.Status == "" {
		input.Status = string(model.MaintenancePlanned)
	}
	if err := validate.MaintenanceInput(input); err != nil {
		return model.MaintenanceTask{}, err
	}
	belongs, err := a.SensorBelongsToSpring(ctx, input.SensorID, input.SpringID)
	if err != nil {
		return model.MaintenanceTask{}, err
	}
	if !belongs {
		return model.MaintenanceTask{}, fmt.Errorf("sensor %d does not belong to spring %d", input.SensorID, input.SpringID)
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	task, err := a.store.CreateMaintenance(dbctx, input, a.Now())
	if err != nil {
		return model.MaintenanceTask{}, a.dbError("create maintenance", err)
	}
	return task, nil
}

func (a *App) GetMaintenance(ctx context.Context, id model.ID) (model.MaintenanceTask, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	task, err := a.store.GetMaintenance(dbctx, id)
	if err != nil {
		return model.MaintenanceTask{}, a.dbError("get maintenance", err)
	}
	return task, nil
}

func (a *App) ListMaintenance(ctx context.Context, filter model.ListFilter) (model.APIList[model.MaintenanceTask], error) {
	filter = filter.Normalized()
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	items, err := a.store.ListMaintenance(dbctx, filter)
	if err != nil {
		return model.APIList[model.MaintenanceTask]{}, a.dbError("list maintenance", err)
	}
	return page(items, filter), nil
}

func (a *App) CompleteMaintenance(ctx context.Context, id model.ID) (model.MaintenanceTask, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	task, err := a.store.CompleteMaintenance(dbctx, id, a.Now())
	if err != nil {
		return model.MaintenanceTask{}, a.dbError("complete maintenance", err)
	}
	return task, nil
}

func (a *App) MarkDueMaintenance(ctx context.Context) (int, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	count, err := a.store.MarkDueMaintenance(dbctx, a.Now())
	if err != nil {
		return 0, a.dbError("mark due maintenance", err)
	}
	return count, nil
}
