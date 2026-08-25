package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

type BatchIngestResult struct {
	Accepted []model.Reading `json:"accepted"`
	Rejected []IngestFailure `json:"rejected"`
}

type IngestFailure struct {
	Index  int    `json:"index"`
	Reason string `json:"reason"`
}

func (a *App) IngestReadings(ctx context.Context, inputs []model.ReadingInput) BatchIngestResult {
	result := BatchIngestResult{Accepted: make([]model.Reading, 0), Rejected: make([]IngestFailure, 0)}
	if len(inputs) > 100 {
		result.Rejected = append(result.Rejected, IngestFailure{Index: -1, Reason: "batch cannot contain more than 100 readings"})
		return result
	}
	for index, input := range inputs {
		if err := ctx.Err(); err != nil {
			result.Rejected = append(result.Rejected, IngestFailure{Index: index, Reason: fmt.Sprintf("batch canceled: %v", err)})
			break
		}
		reading, err := a.IngestReading(ctx, input)
		if err != nil {
			result.Rejected = append(result.Rejected, IngestFailure{Index: index, Reason: err.Error()})
			continue
		}
		result.Accepted = append(result.Accepted, reading)
	}
	return result
}

func (a *App) ValidateReadingBatch(inputs []model.ReadingInput) []IngestFailure {
	failures := make([]IngestFailure, 0)
	for index, input := range inputs {
		if err := validate.ReadingInput(input); err != nil {
			failures = append(failures, IngestFailure{Index: index, Reason: err.Error()})
		}
	}
	return failures
}

func (a *App) AddSamples(ctx context.Context, batchID model.ID, inputs []model.SampleInput) ([]model.Sample, []IngestFailure) {
	items := make([]model.Sample, 0, len(inputs))
	failures := make([]IngestFailure, 0)
	if len(inputs) > 100 {
		failures = append(failures, IngestFailure{Index: -1, Reason: "sample batch cannot contain more than 100 samples"})
		return items, failures
	}
	for index, input := range inputs {
		if err := ctx.Err(); err != nil {
			failures = append(failures, IngestFailure{Index: index, Reason: fmt.Sprintf("batch canceled: %v", err)})
			break
		}
		sample, err := a.AddSample(ctx, batchID, input)
		if err != nil {
			failures = append(failures, IngestFailure{Index: index, Reason: err.Error()})
			continue
		}
		items = append(items, sample)
	}
	return items, failures
}

func (a *App) RequireCompleteBatch(result BatchIngestResult) error {
	if len(result.Rejected) > 0 {
		return fmt.Errorf("batch accepted %d item(s) and rejected %d item(s)", len(result.Accepted), len(result.Rejected))
	}
	return nil
}
