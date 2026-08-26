package engine

import "math"

func lift(previous, baseline, value, ratio float64) bool {
	threshold := math.Max(math.Abs(baseline)*ratio, 0.01)
	return value-previous > threshold || value-baseline > threshold
}

func droppedFromPeak(peak, value, ratio float64) bool {
	if peak == 0 {
		return false
	}
	return (peak-value)/math.Abs(peak) >= ratio
}

func criticalLevel(baseline, peak float64) bool {
	if baseline == 0 {
		return math.Abs(peak) >= 10
	}
	return math.Abs((peak-baseline)/baseline) >= 0.5
}

func warningLevel(baseline, peak float64) bool {
	if baseline == 0 {
		return math.Abs(peak) >= 2
	}
	return math.Abs((peak-baseline)/baseline) >= 0.25
}
