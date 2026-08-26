package clock

import "time"

type Window struct {
	From time.Time
	To   time.Time
}

func LastHours(c Clock, hours int) Window {
	if hours <= 0 {
		hours = 24
	}
	to := c.Now()
	return Window{From: to.Add(-time.Duration(hours) * time.Hour), To: to}
}

func LastDays(c Clock, days int) Window {
	if days <= 0 {
		days = 1
	}
	return LastHours(c, days*24)
}

func CurrentDay(c Clock) Window {
	now := c.Now()
	from := StartOfDay(now)
	return Window{From: from, To: from.Add(24*time.Hour - time.Nanosecond)}
}

func PreviousDay(c Clock) Window {
	today := StartOfDay(c.Now())
	return Window{From: today.Add(-24 * time.Hour), To: today.Add(-time.Nanosecond)}
}

func (w Window) Valid() bool {
	return !w.From.IsZero() && !w.To.IsZero() && !w.From.After(w.To)
}

func (w Window) Contains(at time.Time) bool {
	return w.Valid() && !at.Before(w.From) && !at.After(w.To)
}

func (w Window) Duration() time.Duration {
	if !w.Valid() {
		return 0
	}
	return w.To.Sub(w.From)
}

func (w Window) Limit(maximum time.Duration) Window {
	if maximum <= 0 || w.Duration() <= maximum {
		return w
	}
	return Window{From: w.To.Add(-maximum), To: w.To}
}

func Merge(left, right Window) Window {
	if !left.Valid() {
		return right
	}
	if !right.Valid() {
		return left
	}
	from := left.From
	if right.From.Before(from) {
		from = right.From
	}
	to := left.To
	if right.To.After(to) {
		to = right.To
	}
	return Window{From: from, To: to}
}
