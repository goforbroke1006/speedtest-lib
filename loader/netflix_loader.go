package loader

func NewNetflixLoader() *netflixLoader {
	return &netflixLoader{}
}

var (
	_ NetworkLoader = &netflixLoader{}
)

type netflixLoader struct {
}

func (n netflixLoader) LoadConfig() error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (n netflixLoader) DownloadSink() <-chan float64 {
	//TODO implement me
	//panic("implement me")
	bytesPerSecondSink := make(chan float64)
	go func() {
		bytesPerSecondSink <- 0.1
		close(bytesPerSecondSink)
	}()
	return bytesPerSecondSink
}

func (n netflixLoader) UploadSink() <-chan float64 {
	//TODO implement me
	//panic("implement me")
	bytesPerSecondSink := make(chan float64)
	go func() {
		bytesPerSecondSink <- 0.1
		close(bytesPerSecondSink)
	}()
	return bytesPerSecondSink
}
