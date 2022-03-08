package domain

// NetworkLoader runs a bunch of request to speed-test servers
type NetworkLoader interface {
	LoadConfig() error

	// DownloadSink returns bits-per-seconds updates channel
	DownloadSink() (bits <-chan float64, err error)

	// UploadSink returns bits-per-seconds updates channel
	UploadSink() (bits <-chan float64, err error)
}
