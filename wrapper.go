package speedtest_lib

import (
	"context"
	"fmt"
)

type ProviderKind string

const (
	ProviderKindOokla   = ProviderKind("ookla")
	ProviderKindNetflix = ProviderKind("netflix")
)

func New() *wrapper {
	return &wrapper{
		loaderOokla:   newOoklaLoader(),
		loaderNetflix: newNetflixLoader(),
	}
}

type wrapper struct {
	loaderOokla   NetworkLoader
	loaderNetflix NetworkLoader
}

// DoRequest runs specific loader and collect data
func (w wrapper) DoRequest(ctx context.Context, kind ProviderKind) (download float64, upload float64, err error) {
	var nl NetworkLoader

	switch kind {
	case ProviderKindOokla:
		nl = w.loaderOokla
	case ProviderKindNetflix:
		nl = w.loaderNetflix
	default:
		return 0, 0, fmt.Errorf("unexpected provider: %s", kind)
	}

	if err = nl.LoadConfig(); err != nil {
		return download, upload, err
	}

	dls := nl.DownloadSink()
LoopD:
	for {
		select {
		case <-ctx.Done():
			return download, upload, err
		case c, ok := <-dls:
			if !ok {
				break LoopD
			}
			download = c
		}

	}

	uls := nl.UploadSink()
LoopU:
	for {
		select {
		case <-ctx.Done():
			return download, upload, err
		case c, ok := <-uls:
			if !ok {
				break LoopU
			}
			upload = c
		}

	}

	return download, upload, err
}
