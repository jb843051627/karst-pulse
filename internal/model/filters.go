package model

import "time"

func NewWindow(now time.Time) (time.Time, time.Time) {
	return now.Add(-24 * time.Hour), now
}

func (f ListFilter) HasSpring() bool {
	return f.SpringID > 0
}

func (f ListFilter) HasSensor() bool {
	return f.SensorID > 0
}

func (f ListFilter) HasStatus() bool {
	return f.Status != ""
}

func (f ListFilter) QualityCondition() string {
	return ""
}
