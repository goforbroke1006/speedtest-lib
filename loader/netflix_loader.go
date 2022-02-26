package loader

import (
	"github.com/ddo/go-fast"
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
	_ NetworkLoader = &netflixLoader{}
)

type netflixLoader struct {
	fastCom          *fast.Fast
	urls             []string
	downloadKbpsChan chan float64
	uploadKbpsChan   chan float64
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
	return n.universalSink()
}

func (n netflixLoader) UploadSink() <-chan float64 {
	return n.universalSink()
}

func (n netflixLoader) universalSink() <-chan float64 {
	bytesPerSecondSink := make(chan float64)

	KbpsChan := make(chan float64)

	go func() {
		for kbps := range KbpsChan {
			bytesPerSecondSink <- kbps * 1024
		}
		close(bytesPerSecondSink)
	}()

	go func() {
		if err := n.fastCom.Measure(n.urls, KbpsChan); err != nil {
			close(KbpsChan)
			close(bytesPerSecondSink)
		}
	}()

	return bytesPerSecondSink
}
