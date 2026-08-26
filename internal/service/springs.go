package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (a *App) CreateSpring(ctx context.Context, input model.SpringInput) (model.Spring, error) {
	input = input.Normalize()
	if err := validate.SpringInput(input); err != nil {
		return model.Spring{}, err
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	spring, err := a.store.CreateSpring(dbctx, input, a.Now())
	if err != nil {
		return model.Spring{}, a.dbError("create spring", err)
	}
	return spring, nil
}

func (a *App) GetSpring(ctx context.Context, id model.ID) (model.Spring, error) {
	if id <= 0 {
		return model.Spring{}, fmt.Errorf("spring id must be positive")
	}
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	spring, err := a.store.GetSpring(dbctx, id)
	if err != nil {
		return model.Spring{}, a.dbError("get spring", err)
	}
	return spring, nil
}

func (a *App) ListSprings(ctx context.Context, filter model.ListFilter) (model.APIList[model.Spring], error) {
	filter = filter.Normalized()
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	items, err := a.store.ListSprings(dbctx, filter)
	if err != nil {
		return model.APIList[model.Spring]{}, a.dbError("list springs", err)
	}
	return page(items, filter), nil
}

func (a *App) SpringExists(ctx context.Context, id model.ID) (bool, error) {
	dbctx, cancel := a.withDBTimeout(ctx)
	defer cancel()
	exists, err := a.store.SpringExists(dbctx, id)
	if err != nil {
		return false, a.dbError("check spring", err)
	}
	return exists, nil
}
