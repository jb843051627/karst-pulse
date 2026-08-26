package engine

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/clock"
)

func (r *Runtime) scheduler() {
	defer r.wg.Done()
	for clock.Wait(r.ctx, r.config.MaintenanceTick) {
		if _, err := r.app.MarkDueMaintenance(r.ctx); err != nil && !isCanceled(r.ctx) {
			r.app.Metrics().Inc("worker_failures_total")
		}
	}
}

func isCanceled(ctx context.Context) bool {
	return ctx.Err() != nil
}

func SchedulerDescription(config Config) string {
	return fmt.Sprintf("maintenance scheduler every %s", config.normalized().MaintenanceTick)
}
