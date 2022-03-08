package speedtest_lib

import (
	"github.com/ddo/go-fast"
	"github.com/goforbroke1006/speedtest-lib/domain"

	"github.com/goforbroke1006/speedtest-lib/pkg/content"
)

func newNetflixLoader() *netflixLoader {
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
	fastCom domain.NetflixFastClient
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

func (n netflixLoader) DownloadSink() (<-chan float64, error) {
	bitsPerSecondSink := make(chan float64)

	KbpsChan := make(chan float64)

	go func() {
		for kbps := range KbpsChan {
			bits := content.DataLen(kbps * content.KiloBit).Bits()
			bitsPerSecondSink <- float64(bits)
		}
		close(bitsPerSecondSink)
	}()

	go func() {
		_ = n.fastCom.Measure(n.urls, KbpsChan)
	}()

	return bitsPerSecondSink, nil
}

func (n netflixLoader) UploadSink() (<-chan float64, error) {
	// TODO: stub because fast.com does not provide upload speed check
	bitsPerSecondSink := make(chan float64)
	go func() {
		bitsPerSecondSink <- 0.0
		close(bitsPerSecondSink)
	}()
	return bitsPerSecondSink, nil
}
