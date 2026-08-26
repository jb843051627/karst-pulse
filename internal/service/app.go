package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/clock"
	"github.com/karst-pulse/karst-pulse/internal/metrics"
	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/store"
)

type App struct {
	store   *store.Store
	clock   clock.Clock
	metrics *metrics.Registry
}

func New(database *store.Store, applicationClock clock.Clock, registry *metrics.Registry) *App {
	if applicationClock == nil {
		applicationClock = clock.Real{}
	}
	if registry == nil {
		registry = metrics.New()
	}
	return &App{store: database, clock: applicationClock, metrics: registry}
}

func (a *App) Store() *store.Store {
	return a.store
}

func (a *App) Metrics() *metrics.Registry {
	return a.metrics
}

func (a *App) Now() time.Time {
	return a.clock.Now().UTC()
}

func (a *App) withDBTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func (a *App) dbError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", operation, err)
}

func isMissing(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func page[T any](items []T, filter model.ListFilter) model.APIList[T] {
	if items == nil {
		items = make([]T, 0)
	}
	cloned := make([]T, len(items))
	copy(cloned, items)
	items = cloned
	return model.APIList[T]{Items: items, Page: model.PageInfo{Limit: filter.Limit, Offset: filter.Offset, Count: len(items)}}
}
