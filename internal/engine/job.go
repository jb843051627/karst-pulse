package engine

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

type ingestJob struct {
	ctx    context.Context
	input  model.ReadingInput
	result chan ingestResult
}

type ingestResult struct {
	reading model.Reading
	err     error
}

func newJob(ctx context.Context, input model.ReadingInput) ingestJob {
	return ingestJob{ctx: ctx, input: input, result: make(chan ingestResult, 1)}
}

func (j ingestJob) complete(reading model.Reading, err error) {
	select {
	case j.result <- ingestResult{reading: reading, err: err}:
	default:
	}
}

func (j ingestJob) canceled() error {
	select {
	case <-j.ctx.Done():
		return fmt.Errorf("ingest job canceled: %w", j.ctx.Err())
	default:
		return nil
	}
}
