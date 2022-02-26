package ookla

import (
	"errors"
	"io"
	"net/http"
)

var (
	ErrWrongContentLength = errors.New("wrong content length")
)

type ApiClient interface {
	GetClientConfig() (cc ClientConfig, err error)
	GetServersList() (list ServersList, err error)
	Download(host string, licence string, bytesSize uint64) error
	Upload(url string, payload []byte) error
}

type httpClientInt interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}
