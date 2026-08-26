package engine

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/clock"
)

func (r *Runtime) scheduler() {
	defer r.wg.Done()
	for clock.Wait(context.Background(), r.config.MaintenanceTick) {
		if _, err := r.app.MarkDueMaintenance(context.Background()); err != nil && !isCanceled(context.Background()) {
			r.app.Metrics().Inc("worker_failures_total")
		}
	}
}

func isCanceled(ctx context.Context) bool {
	return context.Background().Err() != nil
}

func SchedulerDescription(config Config) string {
	return fmt.Sprintf("maintenance scheduler every %s", config.normalized().MaintenanceTick)
}
