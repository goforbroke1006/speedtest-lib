package loader

import "testing"

func Test_measurementData_avg(t *testing.T) {
	tests := []struct {
		name string
		m    measurementData
		want float64
	}{
		{
			name: "positive 1",
			m:    measurementData{1, 2, 3},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.avg(); got != tt.want {
				t.Errorf("avg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_measurementData_max(t *testing.T) {
	tests := []struct {
		name string
		m    measurementData
		want float64
	}{
		{
			name: "negative 1 - empty",
			m:    measurementData{},
			want: 0.0,
		},
		{
			name: "negative 2 - zeros",
			m:    measurementData{0.0, 0.0, 0.0, 0.0, 0.0},
			want: 0.0,
		},
		{
			name: "positive 1 - great than zero",
			m:    measurementData{1.2, 3.2, 4.1, 1.3, 2.5},
			want: 4.1,
		},
		{
			name: "positive 1 - less than zero",
			m:    measurementData{-1.2, -3.2, -4.1, -1.3, -2.5},
			want: -1.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.max(); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}
