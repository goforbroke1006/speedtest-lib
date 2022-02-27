package domain

// NetworkLoader runs a bunch of request to speed-test servers
type NetworkLoader interface {
	LoadConfig() error

	// DownloadSink returns bytes-per-seconds updates channel
	DownloadSink() (bytesPerSecondSink <-chan float64)

	// UploadSink returns bytes-per-seconds updates channel
	UploadSink() (bytesPerSecondSink <-chan float64)
}
