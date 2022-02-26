package loader

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/goforbroke1006/speedtest-lib/external/ookla"
)

func NewOoklaLoader() *ooklaLoader {
	client := ookla.NewOoklaSpeedTestClient()

	return &ooklaLoader{
		client:     client,
		servers:    nil,
		licenseKey: "",
		//payloadSizes: []int{
		//	1 * 1024 * 1024,  // 1 mega bytes
		//	2 * 1024 * 1024,  // 2 mega bytes
		//	5 * 1024 * 1024,  // 5 mega bytes
		//	10 * 1024 * 1024, // 10 mega bytes
		//	25 * 1024 * 1024, // 25 mega bytes
		//},
		payloadSize: 25 * 1024 * 1024, // 25 mega bytes
	}
}

var (
	_ NetworkLoader = &ooklaLoader{}
)

type ooklaLoader struct {
	client  ookla.ApiClient
	servers []struct {
		host      string
		uploadUrl string
	}
	licenseKey string
	//payloadSizes []int
	payloadSize uint

	downloadThreadTotal   uint
	downloadThreadPerLink uint
	uploadThreadTotal     uint
	uploadThreadPerLink   uint
}

func (o *ooklaLoader) LoadConfig() error {
	cc, err := o.client.GetClientConfig()
	if err != nil {
		return err
	}

	o.downloadThreadTotal = 1
	o.downloadThreadPerLink = cc.Download.ThreadsPerUrl
	o.uploadThreadTotal = cc.Upload.Threads
	o.uploadThreadPerLink = cc.Upload.ThreadsPerUrl

	o.licenseKey = cc.LicenseKey

	serversList, err := o.client.GetServersList()
	if err != nil {
		return err
	}
	var servers []struct {
		host      string
		uploadUrl string
	}

	for _, s := range serversList {
		servers = append(servers, struct {
			host      string
			uploadUrl string
		}{host: s.Host, uploadUrl: s.Url})
	}
	o.servers = servers

	return nil
}

func (o ooklaLoader) DownloadSink() <-chan float64 {
	bytesPerSecondSink := make(chan float64)
	measurements := measurementData{}

	go func() {
		totalWG := sync.WaitGroup{}
		maxThreadsSem := make(chan struct{}, o.downloadThreadTotal)
		for _, server := range o.servers {

			maxThreadsSem <- struct{}{} // reserve
			totalWG.Add(1)

			go func(host string) {
				defer func() {
					<-maxThreadsSem // release
					totalWG.Done()
				}()

				perLinkWG := sync.WaitGroup{}
				finished := uint(0)
				start := time.Now()
				for i := uint(0); i < o.downloadThreadPerLink; i++ {
					perLinkWG.Add(1)
					go func() {
						defer perLinkWG.Done()
						if err := o.client.Download(host, o.licenseKey, uint64(o.payloadSize)); err != nil {
							fmt.Println(err.Error())
							return
						}
						finished++
					}()
				}
				perLinkWG.Wait()
				speed := float64(o.payloadSize*finished) / time.Since(start).Seconds()
				measurements = append(measurements, speed)

				bytesPerSecondSink <- (measurements.avg() + measurements.max()) / 2

			}(server.host)
		}
		totalWG.Wait()
		close(bytesPerSecondSink)
	}()

	return bytesPerSecondSink
}

func (o ooklaLoader) UploadSink() <-chan float64 {
	bytesPerSecondSink := make(chan float64)
	measurements := measurementData{}

	payload := make([]byte, o.payloadSize)
	rand.Read(payload)

	go func() {
		totalWG := sync.WaitGroup{}
		maxThreadsSem := make(chan struct{}, o.uploadThreadTotal)
		for _, server := range o.servers {

			maxThreadsSem <- struct{}{} // reserve
			totalWG.Add(1)

			go func(url string) {
				defer func() {
					<-maxThreadsSem // release
					totalWG.Done()
				}()

				perLinkWG := sync.WaitGroup{}
				finished := uint(0)
				start := time.Now()
				for i := uint(0); i < o.uploadThreadPerLink; i++ {
					perLinkWG.Add(1)
					go func() {
						defer perLinkWG.Done()
						if err := o.client.Upload(url, payload); err != nil {
							return
						}
						finished++
					}()
				}
				perLinkWG.Wait()
				speed := float64(o.payloadSize*finished) / time.Since(start).Seconds()
				measurements = append(measurements, speed)

				bytesPerSecondSink <- (measurements.avg() + measurements.max()) / 2

			}(server.uploadUrl)
		}
		totalWG.Wait()
		close(bytesPerSecondSink)
	}()

	return bytesPerSecondSink
}
