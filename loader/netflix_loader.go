package loader

import (
	"github.com/ddo/go-fast"

	"github.com/goforbroke1006/speedtest-lib/domain"
)

func NewNetflixLoader() *netflixLoader {
	fastCom := fast.New()
	err := fastCom.Init()
	if err != nil {
		panic(err)
	}

	return &netflixLoader{
		fastCom: fastCom,
	}
}

var (
	_ domain.NetworkLoader = &netflixLoader{}
)

type netflixLoader struct {
	fastCom *fast.Fast
	urls    []string
}

func (n *netflixLoader) LoadConfig() error {
	urls, err := n.fastCom.GetUrls()
	if err != nil {
		return err
	}
	n.urls = urls
	return nil
}

func (n netflixLoader) DownloadSink() <-chan float64 {
	bytesPerSecondSink := make(chan float64)

	KbpsChan := make(chan float64)

	go func() {
		for kbps := range KbpsChan {
			bytesPerSecondSink <- kbps * 1024
		}
		close(bytesPerSecondSink)
	}()

	go func() {
		_ = n.fastCom.Measure(n.urls, KbpsChan)
	}()

	return bytesPerSecondSink
}

func (n netflixLoader) UploadSink() <-chan float64 {
	bytesPerSecondSink := make(chan float64)
	go func() {
		bytesPerSecondSink <- 0.0
		close(bytesPerSecondSink)
	}()
	return bytesPerSecondSink
}
