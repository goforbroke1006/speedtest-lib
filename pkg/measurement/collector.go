package measurement

import "math"

type MetricsCollector []float64

func (m MetricsCollector) Avg() float64 {
	if len(m) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range m {
		sum += v
	}
	return sum / float64(len(m))
}

func (m MetricsCollector) Max() float64 {
	if len(m) == 0 {
		return 0.0
	}
	max := -1 * math.MaxFloat64
	for _, v := range m {
		if v > max {
			max = v
		}
	}
	return max
}
