package validate

import "github.com/karst-pulse/karst-pulse/internal/model"

func BatchInput(input model.BatchInput) error {
	var errors Errors
	if input.SpringID <= 0 {
		errors.Add("spring_id", "must be positive")
	}
	Required(&errors, "batch_code", input.BatchCode)
	Length(&errors, "batch_code", input.BatchCode, 3, 64)
	Required(&errors, "collector", input.Collector)
	TimePresent(&errors, "sampled_at", input.SampledAt)
	if !validBatchStatus(input.Status) {
		errors.Add("status", "must be open, submitted, or archived")
	}
	return errors.Err()
}

func SampleInput(input model.SampleInput) error {
	var errors Errors
	Required(&errors, "parameter", input.Parameter)
	Required(&errors, "unit", input.Unit)
	Finite(&errors, "value", input.Value)
	return errors.Err()
}

func validBatchStatus(value string) bool {
	return value == string(model.BatchOpen) || value == string(model.BatchSubmitted) || value == string(model.BatchArchived)
}
