package regression

import (
    "testing"
    "time"

    "github.com/karst-pulse/karst-pulse/internal/engine"
    "github.com/karst-pulse/karst-pulse/internal/model"
    "github.com/karst-pulse/karst-pulse/internal/transport"
)

func TestBug26_WindowAndPageResultsDoNotAlias(t *testing.T) {
    now := time.Now().UTC()
    aggregator := engine.NewAggregator(time.Hour)
    aggregator.Observe(model.Reading{SpringID: 1, Value: -100, Quality: model.QualityGood, ObservedAt: now.Add(-2 * time.Hour)})
    aggregator.Observe(model.Reading{SpringID: 1, Value: 5, Quality: model.QualityGood, ObservedAt: now.Add(-time.Minute)})
    result := aggregator.Aggregate(1, now)
    if result.Minimum != 5 || result.Maximum != 5 {
        t.Fatalf("window aggregate included stale value: %#v", result)
    }
    source := []int{1, 2}
    page := transport.ApplyPage(source, model.PageInfo{Limit: 2})
    page.Items[0] = 99
    if source[0] != 1 {
        t.Fatalf("page result aliases source slice: %#v", source)
    }
}
