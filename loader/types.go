package loader

import "math"

type measurementData []float64

func (m measurementData) avg() float64 {
	if len(m) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range m {
		sum += v
	}
	return sum / float64(len(m))
}

func (m measurementData) max() float64 {
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
