package upgrader

import (
	"fmt"
	"sync"
	"time"

	"github.com/goforbroke1006/speedtest-lib/loader"
)

func NewUpgrader(nl loader.NetworkLoader, interval time.Duration) *dataUpgrader {
	return &dataUpgrader{
		nl:          nl,
		interval:    interval,
		dlSpeedMbps: 0.0,
		ulSpeedMbps: 0.0,
	}
}

type dataUpgrader struct {
	nl       loader.NetworkLoader
	interval time.Duration

	dlSpeedMbps    float64
	ulSpeedMbps    float64
	measurementsMX sync.RWMutex
}

// Run loads download/upload speed summary asynchronously
func (u *dataUpgrader) Run() {
	for {
		start := time.Now()

		func() {
			err := u.nl.LoadConfig()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			{
				sink := u.nl.DownloadSink()
				for {
					download, opened := <-sink
					if !opened {
						break
					}
					u.measurementsMX.Lock()
					u.dlSpeedMbps = download
					u.measurementsMX.Unlock()
				}
			}

			{
				sink := u.nl.UploadSink()
				for {
					upload, opened := <-sink
					if !opened {
						break
					}
					u.measurementsMX.Lock()
					u.ulSpeedMbps = upload
					u.measurementsMX.Unlock()
				}
			}
		}()

		spend := time.Since(start)
		time.Sleep(u.interval - spend)
	}
}

func (u *dataUpgrader) GetDLSpeedMbps() float64 {
	u.measurementsMX.RLock()
	defer u.measurementsMX.RUnlock()

	return u.dlSpeedMbps
}

func (u *dataUpgrader) GetULSpeedMbps() float64 {
	u.measurementsMX.RLock()
	defer u.measurementsMX.RUnlock()

	return u.ulSpeedMbps
}

func (u *dataUpgrader) IsReady() bool {
	u.measurementsMX.RLock()
	defer u.measurementsMX.RUnlock()

	return u.dlSpeedMbps > 0.0 && u.ulSpeedMbps > 0.0
}
