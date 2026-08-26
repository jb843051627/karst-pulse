package regression

import (
    "context"
    "errors"
    "path/filepath"
    "testing"
    "time"

    "github.com/karst-pulse/karst-pulse/internal/clock"
    "github.com/karst-pulse/karst-pulse/internal/metrics"
    "github.com/karst-pulse/karst-pulse/internal/model"
    "github.com/karst-pulse/karst-pulse/internal/service"
    "github.com/karst-pulse/karst-pulse/internal/store"
)

func TestBug30_AlertRelationsAndCancellationStayIsolated(t *testing.T) {
    sensorID := model.ID(7)
    input := model.AlertInput{SpringID: 1, SensorID: &sensorID, Message: "读数异常", TriggeredAt: time.Now().UTC()}
    alert := input.Normalize(time.Now().UTC())
    sensorID = 99
    if alert.SensorID == nil || *alert.SensorID != 7 {
        t.Fatalf("alert relation aliases caller pointer: %#v", alert.SensorID)
    }
    database, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "karst.db"))
    if err != nil {
        t.Fatal(err)
    }
    defer database.Close()
    app := service.New(database, clock.Real{}, metrics.New())
    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    _, err = app.CreateAlert(ctx, model.Alert{SpringID: 1, Message: "cancel me", Status: model.AlertOpen})
    if !errors.Is(err, context.Canceled) {
        t.Fatalf("cancellation identity lost: %v", err)
    }
}
