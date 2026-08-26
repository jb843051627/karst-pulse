package clock

import "time"

type Real struct{}

func (Real) Now() time.Time {
	return time.Now().UTC()
}

func StartOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func EndOfDay(value time.Time) time.Time {
	return StartOfDay(value).Add(24*time.Hour - time.Nanosecond)
}
