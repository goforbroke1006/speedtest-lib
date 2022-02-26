package test_speed

type Upgrader interface {
	GetDLSpeedMbps() float64
	GetULSpeedMbps() float64
}
