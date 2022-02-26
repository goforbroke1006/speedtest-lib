package domain

type NetworkLoader interface {
	LoadConfig() error

	// DownloadSink returns bytes-per-seconds updates channel
	DownloadSink() (bytesPerSecondSink <-chan float64)

	// UploadSink returns bytes-per-seconds updates channel
	UploadSink() (bytesPerSecondSink <-chan float64)
}
