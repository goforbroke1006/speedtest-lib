package loader

type NetworkLoader interface {
	LoadServersList() (uint, error)
	Download() (bytesPerSecond float64, err error)
	Upload() (bytesPerSecond float64, err error)
}
