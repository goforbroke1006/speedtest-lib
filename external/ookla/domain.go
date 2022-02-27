package ookla

import (
	"errors"
	"io"
	"net/http"
)

var (
	ErrWrongContentLength = errors.New("wrong content length")
)

// ApiClient wraps requests to https://www.speedtest.net/ API
type ApiClient interface {

	// GetClientConfig loads client configuration and loading plan from external API
	GetClientConfig() (cc ClientConfig, err error)

	// GetServersList loads servers list from external API
	GetServersList() (list ServersList, err error)

	// Download allows asking server for downloading content with certain length
	Download(host string, licence string, bytesSize uint64) error

	// Upload allow sending content to server
	Upload(url string, payload []byte) error
}

// httpClientInt is custom http.Client interface
// It allows writing test with mocks
type httpClientInt interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}
