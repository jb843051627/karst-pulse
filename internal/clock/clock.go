package clock

import "time"

type Clock interface {
	Now() time.Time
}

func Today(c Clock) time.Time {
	now := c.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func Since(c Clock, duration time.Duration) time.Time {
	return c.Now().Add(-duration)
}
