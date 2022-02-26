package loader

type measurementData []float64

func (m measurementData) arg() float64 {
	if len(m) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range m {
		sum += v
	}
	return sum / float64(len(m))
}
