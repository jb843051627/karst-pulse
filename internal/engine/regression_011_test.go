package engine

import (
    "sync"
    "testing"
    "time"

    "github.com/karst-pulse/karst-pulse/internal/model"
)

func TestBug11_AggregatorProtectsSharedSeries(t *testing.T) {
    aggregator := NewAggregator(time.Hour)
    now := time.Now().UTC()
    var wg sync.WaitGroup
    for worker := 0; worker < 20; worker++ {
        wg.Add(1)
        go func(worker int) {
            defer wg.Done()
            for index := 0; index < 100; index++ {
                aggregator.Observe(model.Reading{SpringID: 1, SensorID: 1, Value: float64(worker + index), Quality: model.QualityGood, ObservedAt: now})
                _ = aggregator.Aggregate(1, now)
                _ = aggregator.Series(1, now.Add(-time.Minute), now.Add(time.Minute))
            }
        }(worker)
    }
    wg.Wait()
}
