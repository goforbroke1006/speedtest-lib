package speedtest_lib

import (
	"github.com/showwin/speedtest-go/speedtest"

	"github.com/goforbroke1006/speedtest-lib/pkg/measurement"
)

func newOoklaLoader() *ooklaLoaderDefault {
	return &ooklaLoaderDefault{}
}

var (
	_ NetworkLoader = &ooklaLoaderDefault{}
)

type ooklaLoaderDefault struct {
	targets speedtest.Servers
}

func (o *ooklaLoaderDefault) LoadConfig() error {
	user, err := speedtest.FetchUserInfo()
	if err != nil {
		return err
	}
	client := speedtest.New()
	serverList, err := client.FetchServers(user)
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

func (o ooklaLoaderDefault) DownloadSink() <-chan float64 {
	bytesPerSecondSink := make(chan float64)
	go func() {
		measurements := measurement.MetricsCollector{}
		for _, s := range o.targets {
			if err := s.PingTest(); err != nil {
				continue
			}
			if err := s.DownloadTest(false); err != nil {
				continue
			}
			measurements = append(measurements, s.DLSpeed)
			//fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)
			bytesPerSecondSink <- measurements.Avg()
		}
		close(bytesPerSecondSink)
	}()
	return bytesPerSecondSink
}

func (o ooklaLoaderDefault) UploadSink() <-chan float64 {
	bytesPerSecondSink := make(chan float64)
	go func() {
		measurements := measurement.MetricsCollector{}
		for _, s := range o.targets {
			if err := s.PingTest(); err != nil {
				continue
			}
			if err := s.UploadTest(false); err != nil {
				continue
			}
			measurements = append(measurements, s.ULSpeed)
			//fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)
			bytesPerSecondSink <- measurements.Avg()
		}
		close(bytesPerSecondSink)
	}()
	return bytesPerSecondSink
}
