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


func TestBug08_MissingSamplingBatchErrorKeepsIdentity(t *testing.T) {
	db, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "karst.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	app := service.New(db, clock.Real{}, metrics.New())
	_, err = app.AddSample(context.Background(), 99999, model.SampleInput{Parameter: "conductivity", Value: 1, Unit: "mS/cm"})
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("missing sql.ErrNoRows identity: %v", err)
	}
}
