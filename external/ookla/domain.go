package ookla

import "errors"

var (
	ErrWrongContentLength = errors.New("wrong content length")
)

type ApiClient interface {
	GetClientConfig() (cc ClientConfig, err error)
	GetServersList() (list ServersList, err error)
	Download(host string, licence string, bytesSize uint64) error
	Upload(url string, payload []byte) error
}
