package engine

import (
    "testing"
    "time"

    "github.com/karst-pulse/karst-pulse/internal/model"
)

func TestBug29_TrimmedSeriesDoesNotAliasState(t *testing.T) {
    now := time.Now().UTC()
    original := []model.Reading{
        {SpringID: 1, Value: 1, ObservedAt: now.Add(-time.Hour)},
        {SpringID: 1, Value: 2, ObservedAt: now},
    }
    trimmed := trimReadings(original, now.Add(-30*time.Minute))
    trimmed[0].Value = 99
    if original[1].Value == 99 {
        t.Fatal("trimmed slice aliases original readings")
    }
    aggregator := NewAggregator(time.Hour)
    aggregator.readings[1] = original
    series := aggregator.Series(1, now.Add(-2*time.Hour), now.Add(time.Hour))
    series[0].Value = 77
    if aggregator.readings[1][0].Value == 77 {
        t.Fatal("trend result aliases aggregator state")
    }
}
