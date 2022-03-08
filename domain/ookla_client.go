package domain

import "github.com/showwin/speedtest-go/speedtest"

type OoklaSpeedTestClient interface {
	FetchUserInfo() (*speedtest.User, error)
	FetchServers(user *speedtest.User) (speedtest.Servers, error)
}
