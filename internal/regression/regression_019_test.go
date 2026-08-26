package regression


import (
    "context"
    "database/sql"
    "errors"
    "path/filepath"
    "strings"
    "testing"
    "time"

    "github.com/karst-pulse/karst-pulse/internal/clock"
    "github.com/karst-pulse/karst-pulse/internal/engine"
    "github.com/karst-pulse/karst-pulse/internal/metrics"
    "github.com/karst-pulse/karst-pulse/internal/model"
    "github.com/karst-pulse/karst-pulse/internal/service"
    "github.com/karst-pulse/karst-pulse/internal/store"
)

var _ = errors.Is
var _ = strings.Contains
var _ = sql.ErrNoRows
var _ = time.Now
var _ = model.ID(0)
var _ = engine.DefaultConfig



func setupApp19(t *testing.T) (*service.App, *store.Store, model.Spring, model.Sensor) {
	t.Helper()
	db, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "karst.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	app := service.New(db, clock.Real{}, metrics.New())
	spring, err := app.CreateSpring(context.Background(), model.SpringInput{
		Code: "SPR-19", Name: "白石泉", Region: "西岭", Aquifer: "灰岩含水层",
		Latitude: 25.1, Longitude: 108.2,
	})
	if err != nil {
		t.Fatal(err)
	}
	sensor, err := app.CreateSensor(context.Background(), model.SensorInput{
		SpringID: spring.ID, SerialNo: "FLOW-19", Kind: string(model.SensorFlow), Unit: "m3/s",
		ThresholdLow: 0, ThresholdHigh: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	return app, db, spring, sensor
}


func TestBug19_ConfirmedPulseTimelineUsesEndTime(t *testing.T) {
	app, _, spring, _ := setupApp19(t)
	start := time.Now().UTC().Add(-time.Hour)
	event, err := app.CreateEvent(context.Background(), model.PulseEvaluation{SpringID: spring.ID, Phase: model.PhaseRising, Severity: model.SeverityWarning, Baseline: 1, PeakValue: 2, At: start, Summary: "rise"})
	if err != nil {
		t.Fatal(err)
	}
	peak := start.Add(10 * time.Minute)
	ended := start.Add(40 * time.Minute)
	event.PeakedAt = &peak
	event.EndedAt = &ended
	event.Phase = model.PhaseConfirmed
	event.UpdatedAt = ended
	if err := app.UpdateEvent(context.Background(), event); err != nil {
		t.Fatal(err)
	}
	timeline, err := app.PulseTimeline(context.Background(), event.ID)
	if err != nil || timeline.Duration != 40*time.Minute || !timeline.Completed {
		t.Fatalf("timeline=%#v err=%v", timeline, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = app.PulseTimeline(ctx, event.ID)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation identity lost: %v", err)
	}
}
