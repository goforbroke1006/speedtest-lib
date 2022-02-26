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
		payloadSizes: []int{
			1 * 1024 * 1024,  // 1 mega bytes
			2 * 1024 * 1024,  // 2 mega bytes
			5 * 1024 * 1024,  // 5 mega bytes
			10 * 1024 * 1024, // 10 mega bytes
			25 * 1024 * 1024, // 25 mega bytes
		},
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
	licenseKey   string
	payloadSizes []int
}

func (o *ooklaLoader) LoadServersList() (uint, error) {
	cc, err := o.client.GetClientConfig()
	if err != nil {
		return 0, err
	}
	o.licenseKey = cc.LicenseKey

	serversList, err := o.client.GetServersList()
	if err != nil {
		return 0, err
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

	return uint(len(servers)), nil
}

func (o ooklaLoader) Download() (bytesPerSecond float64, err error) {
	measurements := measurementData{}

	var wg sync.WaitGroup
	for _, server := range o.servers {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			for _, ps := range o.payloadSizes {
				start := time.Now()
				if err := o.client.Download(host, o.licenseKey, uint64(ps)); err != nil {
					fmt.Println(err.Error())
					return
				}
				spend := time.Since(start).Seconds()
				speed := float64(ps) / spend
				measurements = append(measurements, speed)
			}

		}(server.host)
	}
	wg.Wait()

	return measurements.arg(), nil
}

func (o ooklaLoader) Upload() (bytesPerSecond float64, err error) {
	measurements := measurementData{}

	var wg sync.WaitGroup
	for _, server := range o.servers {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			for _, ps := range o.payloadSizes {
				payload := make([]byte, ps)
				rand.Read(payload)

				start := time.Now()
				if err := o.client.Upload(url, payload); err != nil {
					return
				}
				spend := time.Since(start).Seconds()
				speed := float64(ps) / spend
				measurements = append(measurements, speed)
			}

		}(server.uploadUrl)
	}
	wg.Wait()

	return measurements.arg(), nil
}
