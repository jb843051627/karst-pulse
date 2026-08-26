package engine

import "time"

type Window struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func NewWindow(to time.Time, duration time.Duration) Window {
	if duration <= 0 {
		duration = 24 * time.Hour
	}
	return Window{From: to.Add(-duration), To: to}
}

func (w Window) Valid() bool {
	return !w.From.IsZero() && !w.To.IsZero() && !w.From.After(w.To)
}

func (w Window) Duration() time.Duration {
	if !w.Valid() {
		return 0
	}
	return w.To.Sub(w.From)
}

func (w Window) Contains(at time.Time) bool {
	return w.Valid() && !at.Before(w.From) && !at.After(w.To)
}

func (w Window) Expand(before, after time.Duration) Window {
	return Window{From: w.From.Add(-before), To: w.To.Add(after)}
}

func (w Window) Shift(duration time.Duration) Window {
	return Window{From: w.From.Add(duration), To: w.To.Add(duration)}
}

func SplitWindow(window Window, parts int) []Window {
	if !window.Valid() || parts <= 0 {
		return []Window{}
	}
	step := window.Duration() / time.Duration(parts)
	if step <= 0 {
		return []Window{window}
	}
	items := make([]Window, 0, parts)
	start := window.From
	for index := 0; index < parts; index++ {
		end := start.Add(step)
		if index == parts-1 || end.After(window.To) {
			end = window.To
		}
		items = append(items, Window{From: start, To: end})
		start = end
	}
	return items
}

func ClampTime(value time.Time, window Window) time.Time {
	if value.Before(window.From) {
		return window.From
	}
	if value.After(window.To) {
		return window.To
	}
	return value
}

func Freshness(now, observed time.Time, threshold time.Duration) string {
	if observed.IsZero() {
		return "missing"
	}
	age := now.Sub(observed)
	if age < 0 {
		return "future"
	}
	if threshold <= 0 {
		threshold = 6 * time.Hour
	}
	if age > threshold {
		return "stale"
	}
	return "fresh"
}
