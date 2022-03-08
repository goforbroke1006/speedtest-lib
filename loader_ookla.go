package speedtest_lib

import (
	"github.com/showwin/speedtest-go/speedtest"

	"github.com/goforbroke1006/speedtest-lib/domain"
	"github.com/goforbroke1006/speedtest-lib/pkg/content"
	"github.com/goforbroke1006/speedtest-lib/pkg/measurement"
)

func newOoklaLoader() *ooklaLoaderDefault {
	return &ooklaLoaderDefault{
		client: speedtest.New(),
	}
}

var (
	_ domain.NetworkLoader = &ooklaLoaderDefault{}
)

type ooklaLoaderDefault struct {
	client  domain.OoklaSpeedTestClient
	targets speedtest.Servers
}

func (o *ooklaLoaderDefault) LoadConfig() error {
	user, err := o.client.FetchUserInfo()
	if err != nil {
		return err
	}

	serverList, err := o.client.FetchServers(user)
	if err != nil {
		return err
	}

	targets, err := serverList.FindServer([]int{})
	if err != nil {
		return err
	}

	o.targets = targets
	return nil
}

func (o ooklaLoaderDefault) DownloadSink() (<-chan float64, error) {
	bitsPerSecondSink := make(chan float64)
	go func() {
		measurements := measurement.MetricsCollector{}
		for _, s := range o.targets {
			if err := s.PingTest(); err != nil {
				continue
			}
			if err := s.DownloadTest(false); err != nil {
				continue
			}
			bits := content.DataLen(s.DLSpeed * content.MegaBit).Bits()
			measurements = append(measurements, float64(bits))
			bitsPerSecondSink <- measurements.Avg()
		}
		close(bitsPerSecondSink)
	}()
	return bitsPerSecondSink, nil
}

func (o ooklaLoaderDefault) UploadSink() (<-chan float64, error) {
	bitsPerSecondSink := make(chan float64)
	go func() {
		measurements := measurement.MetricsCollector{}
		for _, s := range o.targets {
			if err := s.PingTest(); err != nil {
				continue
			}
			if err := s.UploadTest(false); err != nil {
				continue
			}
			bits := content.DataLen(s.ULSpeed * content.MegaBit).Bits()
			measurements = append(measurements, float64(bits))
			bitsPerSecondSink <- measurements.Avg()
		}
		close(bitsPerSecondSink)
	}()
	return bitsPerSecondSink, nil
}
