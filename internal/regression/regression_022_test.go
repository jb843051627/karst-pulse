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



func setupApp22(t *testing.T) (*service.App, *store.Store, model.Spring, model.Sensor) {
	t.Helper()
	db, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "karst.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	app := service.New(db, clock.Real{}, metrics.New())
	spring, err := app.CreateSpring(context.Background(), model.SpringInput{
		Code: "SPR-22", Name: "白石泉", Region: "西岭", Aquifer: "灰岩含水层",
		Latitude: 25.1, Longitude: 108.2,
	})
	if err != nil {
		t.Fatal(err)
	}
	sensor, err := app.CreateSensor(context.Background(), model.SensorInput{
		SpringID: spring.ID, SerialNo: "FLOW-22", Kind: string(model.SensorFlow), Unit: "m3/s",
		ThresholdLow: 0, ThresholdHigh: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	return app, db, spring, sensor
}


func TestBug22_MaintenanceCompletionIsIdempotent(t *testing.T) {
	app, _, spring, sensor := setupApp22(t)
	task, err := app.CreateMaintenance(context.Background(), model.MaintenanceInput{SpringID: spring.ID, SensorID: sensor.ID, Title: "清理探头", DueAt: time.Now().UTC().Add(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := app.CompleteMaintenance(ctx, task.ID); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled completion succeeded or lost identity: %v", err)
	}
	if _, err := app.CompleteMaintenance(context.Background(), task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := app.CompleteMaintenance(context.Background(), task.ID); err == nil {
		t.Fatal("completed maintenance was accepted a second time")
	}
}
