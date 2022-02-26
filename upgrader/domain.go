package upgrader

type Upgrader interface {
	GetDLSpeedMbps() float64
	GetULSpeedMbps() float64
	IsReady() bool
}
