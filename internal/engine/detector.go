package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/service"
)

type Detector struct {
	app    *service.App
	config Config
	cache  *stateCache
	mu     sync.Mutex
}

func NewDetector(app *service.App, config Config) *Detector {
	return &Detector{app: app, config: config.normalized(), cache: newStateCache()}
}

func (d *Detector) Observe(ctx context.Context, reading model.Reading) error {
	if !reading.IsUsable() {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("observe pulse canceled: %w", err)
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	state, exists := d.cache.get(reading.SpringID)
	if !exists || state.observed == 0 {
		state = springState{previous: reading.Value, baseline: reading.Value, peak: reading.Value, observedAt: reading.ObservedAt, observed: 1}
		d.cache.put(reading.SpringID, state)
		return nil
	}
	evaluation := d.evaluate(state, reading)
	state = d.advance(state, reading, evaluation.Phase)
	d.cache.put(reading.SpringID, state)
	if evaluation.Phase == "" {
		return nil
	}
	return d.persistEvaluation(ctx, evaluation)
}

func (d *Detector) evaluate(state springState, reading model.Reading) model.PulseEvaluation {
	evaluation := model.PulseEvaluation{SpringID: reading.SpringID, Baseline: state.baseline, PeakValue: maxFloat(state.peak, reading.Value), At: reading.ObservedAt}
	if !state.active && lift(state.previous, state.baseline, reading.Value, d.config.PulseLift) {
		evaluation.Phase = model.PhaseRising
		evaluation.Severity = d.severity(state.baseline, reading.Value)
		evaluation.Summary = fmt.Sprintf("水文脉冲上升，读数 %.3f，高于基线 %.3f", reading.Value, state.baseline)
		return evaluation
	}
	if state.active && reading.Value >= state.peak {
		evaluation.Phase = model.PhasePeaked
		evaluation.Severity = d.severity(state.baseline, reading.Value)
		evaluation.Summary = fmt.Sprintf("水文脉冲达到峰值 %.3f", reading.Value)
		return evaluation
	}
	if state.active && droppedFromPeak(state.peak, reading.Value, d.config.PulseDrop) {
		evaluation.Phase = model.PhaseFading
		evaluation.Severity = d.severity(state.baseline, state.peak)
		evaluation.Summary = fmt.Sprintf("水文脉冲衰减，当前读数 %.3f", reading.Value)
		if !lift(state.baseline, state.baseline, reading.Value, d.config.PulseLift/2) {
			evaluation.Phase = model.PhaseConfirmed
			evaluation.EndAt = reading.ObservedAt
		}
		return evaluation
	}
	return model.PulseEvaluation{}
}

func (d *Detector) advance(state springState, reading model.Reading, phase model.PulsePhase) springState {
	state.previous = reading.Value
	state.observedAt = reading.ObservedAt
	state.observed++
	if reading.Value > state.peak {
		state.peak = reading.Value
	}
	if phase == model.PhaseRising {
		state.active = true
	}
	if phase == model.PhaseConfirmed {
		state.active = false
		state.baseline = reading.Value
		state.peak = reading.Value
	}
	return state
}

func (d *Detector) severity(baseline, peak float64) model.EventSeverity {
	if criticalLevel(baseline, peak) {
		return model.SeverityCritical
	}
	if warningLevel(baseline, peak) {
		return model.SeverityWarning
	}
	return model.SeverityInfo
}

func maxFloat(left, right float64) float64 {
	if left > right {
		return left
	}
	return right
}

func (d *Detector) persistEvaluation(ctx context.Context, evaluation model.PulseEvaluation) error {
	open, err := d.app.LatestOpenEvent(ctx, evaluation.SpringID)
	if err != nil {
		return fmt.Errorf("load open event: %w", err)
	}
	if open == nil {
		event, createErr := d.app.CreateEvent(ctx, evaluation)
		if createErr != nil {
			return createErr
		}
		if event.NeedsAttention() {
			_, alertErr := d.app.CreateAlert(ctx, model.Alert{SpringID: event.SpringID, EventID: &event.ID, Kind: "pulse", Level: model.AlertLevel(event.Severity), Status: model.AlertOpen, Message: event.Summary, TriggeredAt: event.StartedAt})
			if alertErr != nil {
				return fmt.Errorf("create pulse alert: %w", alertErr)
			}
		}
		return nil
	}
	open.Phase = evaluation.Phase
	open.Severity = evaluation.Severity
	open.PeakValue = maxFloat(open.PeakValue, evaluation.PeakValue)
	open.Summary = evaluation.Summary
	open.UpdatedAt = d.app.Now()
	if evaluation.Phase == model.PhasePeaked {
		now := evaluation.At
		open.PeakedAt = &now
	}
	if evaluation.Phase == model.PhaseConfirmed {
		now := evaluation.At
		open.EndedAt = &now
	}
	return d.app.UpdateEvent(ctx, *open)
}

func (d *Detector) CacheSize() int {
	return d.cache.size()
}
