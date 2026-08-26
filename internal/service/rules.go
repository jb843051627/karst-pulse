package service

import (
	"context"
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

type ReadingAssessment struct {
	ReadingID model.ID         `json:"reading_id"`
	SensorID  model.ID         `json:"sensor_id"`
	Value     float64          `json:"value"`
	State     string           `json:"state"`
	Level     model.AlertLevel `json:"level"`
	Message   string           `json:"message"`
}

func (a *App) AssessReading(ctx context.Context, reading model.Reading) (ReadingAssessment, error) {
	sensor, err := a.GetSensor(ctx, reading.SensorID)
	if err != nil {
		return ReadingAssessment{}, fmt.Errorf("load sensor for assessment: %w", err)
	}
	assessment := ReadingAssessment{ReadingID: reading.ID, SensorID: reading.SensorID, Value: reading.Value, State: "normal", Level: model.AlertInfo}
	if reading.Value < sensor.ThresholdLow {
		assessment.State = "below_range"
		assessment.Level = model.AlertWarning
		assessment.Message = fmt.Sprintf("传感器 %s 读数 %.3f 低于下限 %.3f", sensor.SerialNo, reading.Value, sensor.ThresholdLow)
	}
	if reading.Value > sensor.ThresholdHigh {
		assessment.State = "above_range"
		assessment.Level = model.AlertWarning
		assessment.Message = fmt.Sprintf("传感器 %s 读数 %.3f 高于上限 %.3f", sensor.SerialNo, reading.Value, sensor.ThresholdHigh)
	}
	if sensor.Status != model.SensorOnline {
		assessment.State = "sensor_unavailable"
		assessment.Level = model.AlertWarning
		assessment.Message = fmt.Sprintf("传感器 %s 当前状态为 %s", sensor.SerialNo, sensor.Status)
	}
	return assessment, nil
}

func (a *App) EvaluateAndAlert(ctx context.Context, reading model.Reading) (ReadingAssessment, *model.Alert, error) {
	assessment, err := a.AssessReading(ctx, reading)
	if err != nil {
		return ReadingAssessment{}, nil, err
	}
	if assessment.State == "normal" {
		return assessment, nil, nil
	}
	alert, err := a.CreateAlert(ctx, model.Alert{
		SpringID:    reading.SpringID,
		SensorID:    &reading.SensorID,
		Kind:        "sensor_threshold",
		Level:       assessment.Level,
		Status:      model.AlertOpen,
		Message:     assessment.Message,
		TriggeredAt: reading.ObservedAt,
	})
	if err != nil {
		return ReadingAssessment{}, nil, fmt.Errorf("create sensor threshold alert: %w", err)
	}
	return assessment, &alert, nil
}

func (a *App) ThresholdState(ctx context.Context, sensorID model.ID, value float64) (string, error) {
	sensor, err := a.GetSensor(ctx, sensorID)
	if err != nil {
		return "", err
	}
	if value < sensor.ThresholdLow {
		return "below_range", nil
	}
	if value > sensor.ThresholdHigh {
		return "above_range", nil
	}
	return "normal", nil
}
