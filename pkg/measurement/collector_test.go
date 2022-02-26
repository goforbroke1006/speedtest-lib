package measurement

import "testing"

func Test_measurementData_avg(t *testing.T) {
	tests := []struct {
		name string
		m    MetricsCollector
		want float64
	}{
		{
			name: "positive 1",
			m:    MetricsCollector{1, 2, 3},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Avg(); got != tt.want {
				t.Errorf("Avg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_measurementData_max(t *testing.T) {
	tests := []struct {
		name string
		m    MetricsCollector
		want float64
	}{
		{
			name: "negative 1 - empty",
			m:    MetricsCollector{},
			want: 0.0,
		},
		{
			name: "negative 2 - zeros",
			m:    MetricsCollector{0.0, 0.0, 0.0, 0.0, 0.0},
			want: 0.0,
		},
		{
			name: "positive 1 - great than zero",
			m:    MetricsCollector{1.2, 3.2, 4.1, 1.3, 2.5},
			want: 4.1,
		},
		{
			name: "positive 1 - less than zero",
			m:    MetricsCollector{-1.2, -3.2, -4.1, -1.3, -2.5},
			want: -1.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Max(); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}
