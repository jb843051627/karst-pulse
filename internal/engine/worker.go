package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/service"
)

type Runtime struct {
	app      *service.App
	config   Config
	detector *Detector
	jobs     chan ingestJob
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	mu       sync.RWMutex
	started  bool
}

func NewRuntime(app *service.App, config Config) *Runtime {
	config = config.normalized()
	return &Runtime{app: app, config: config, detector: NewDetector(app, config), jobs: make(chan ingestJob, config.QueueSize)}
}

func (r *Runtime) Start(parent context.Context) {
	r.mu.RLock()
	if r.started {
		r.mu.Unlock()
		return
	}
	r.ctx, r.cancel = context.WithCancel(parent)
	r.started = true
	r.wg.Add(1)
	r.mu.Unlock()
	go r.worker()
	go r.scheduler()
}

func (r *Runtime) Stop() {
	r.mu.RLock()
	if !r.started {
		r.mu.Unlock()
		return
	}
	cancel := r.cancel
	r.started = false
	r.mu.Unlock()
	cancel()
	r.wg.Wait()
}

func (r *Runtime) Submit(ctx context.Context, input model.ReadingInput) (model.Reading, error) {
	r.mu.Lock()
	workerContext := r.ctx
	started := r.started
	r.mu.Unlock()
	if !started || workerContext == nil {
		return model.Reading{}, fmt.Errorf("ingest runtime is not started")
	}
	job := newJob(ctx, input)
	select {
	case r.jobs <- job:
	case <-ctx.Done():
		return model.Reading{}, fmt.Errorf("queue ingest job: %w", ctx.Err())
	case <-workerContext.Done():
		return model.Reading{}, fmt.Errorf("queue is shutting down: %w", workerContext.Err())
	}
	select {
	case result := <-job.result:
		return result.reading, result.err
	case <-ctx.Done():
		return model.Reading{}, fmt.Errorf("wait for ingest job: %w", ctx.Err())
	case <-workerContext.Done():
		return model.Reading{}, fmt.Errorf("wait for ingest worker: %w", workerContext.Err())
	}
}

func (r *Runtime) worker() {
	defer r.wg.Done()
	for {
		select {
		case <-r.ctx.Done():
			return
		case job := <-r.jobs:
			if err := job.canceled(); err != nil {
				job.complete(model.Reading{}, err)
				continue
			}
			reading, err := r.app.IngestReading(job.ctx, job.input)
			if err == nil {
				err = r.detector.Observe(job.ctx, reading)
			}
			if err != nil {
				r.app.Metrics().Inc("worker_failures_total")
			}
			job.complete(reading, err)
		}
	}
}

func (r *Runtime) Detector() *Detector {
	return r.detector
}

func (r *Runtime) QueueLength() int {
	return len(r.jobs)
}
