package upgrader

type Upgrader interface {
	Run()
	GetDLSpeedMbps() float64
	GetULSpeedMbps() float64
	IsReady() bool
}
