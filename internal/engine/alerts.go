package engine

import (
	"fmt"

	"github.com/karst-pulse/karst-pulse/internal/model"
)

func AlertLevelForSeverity(severity model.EventSeverity) model.AlertLevel {
	switch severity {
	case model.SeverityCritical:
		return model.AlertCritical
	case model.SeverityWarning:
		return model.AlertWarning
	default:
		return model.AlertInfo
	}
}

func ShouldAlert(severity model.EventSeverity, repeated bool) bool {
	if severity == model.SeverityCritical {
		return true
	}
	return severity == model.SeverityWarning && !repeated
}

func PulseAlert(event model.PulseEvent, repeated bool) (model.Alert, bool) {
	if !ShouldAlert(event.Severity, repeated) {
		return model.Alert{}, false
	}
	return model.Alert{
		SpringID:    event.SpringID,
		EventID:     &event.ID,
		Kind:        "pulse",
		Level:       AlertLevelForSeverity(event.Severity),
		Status:      model.AlertOpen,
		Message:     event.Summary,
		TriggeredAt: event.StartedAt,
	}, true
}

func SensorAlertMessage(sensor model.Sensor, value float64, above bool) string {
	if above {
		return fmt.Sprintf("传感器 %s 读数 %.3f 高于阈值 %.3f", sensor.SerialNo, value, sensor.ThresholdHigh)
	}
	return fmt.Sprintf("传感器 %s 读数 %.3f 低于阈值 %.3f", sensor.SerialNo, value, sensor.ThresholdLow)
}
