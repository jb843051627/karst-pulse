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



func setupApp25(t *testing.T) (*service.App, *store.Store, model.Spring, model.Sensor) {
	t.Helper()
	db, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "karst.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	app := service.New(db, clock.Real{}, metrics.New())
	spring, err := app.CreateSpring(context.Background(), model.SpringInput{
		Code: "SPR-25", Name: "白石泉", Region: "西岭", Aquifer: "灰岩含水层",
		Latitude: 25.1, Longitude: 108.2,
	})
	if err != nil {
		t.Fatal(err)
	}
	sensor, err := app.CreateSensor(context.Background(), model.SensorInput{
		SpringID: spring.ID, SerialNo: "FLOW-25", Kind: string(model.SensorFlow), Unit: "m3/s",
		ThresholdLow: 0, ThresholdHigh: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	return app, db, spring, sensor
}


func TestBug25_SuccessfulReadingTransactionCommits(t *testing.T) {
	app, _, spring, sensor := setupApp25(t)
	reading, err := app.IngestReading(context.Background(), model.ReadingInput{SpringID: spring.ID, SensorID: sensor.ID, Value: 4.2, ObservedAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := app.GetReading(context.Background(), reading.ID)
	if err != nil || loaded.Value != 4.2 {
		t.Fatalf("committed reading missing: %#v err=%v", loaded, err)
	}
}
