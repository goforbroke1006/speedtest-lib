package ookla

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func NewOoklaSpeedTestClient() *ooklaSpeedTestClient {
	hc := http.Client{
		Timeout: time.Minute,
	}
	return &ooklaSpeedTestClient{
		httpClient: hc,
	}
}

type ooklaSpeedTestClient struct {
	httpClient http.Client
}

func (c ooklaSpeedTestClient) GetClientConfig() (cc ClientConfig, err error) {
	var (
		resp     *http.Response
		respBody []byte
	)

	resp, err = c.httpClient.Get(clientConfig)
	if err != nil {
		return
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = xml.Unmarshal(respBody, &cc)

	return cc, err
}

func (c ooklaSpeedTestClient) GetServersList() (list ServersList, err error) {
	resp, err := c.httpClient.Get(speedTestServersUrl)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &list)

	return list, err
}

func (c ooklaSpeedTestClient) Download(host string, licence string, bytesSize uint64) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://%s/download?nocache=%s&size=%d&guid=%s", host, newUUID.String(), bytesSize, licence)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if uint64(len(respBody)) != bytesSize {
		return ErrWrongContentLength
	}

	return nil
}

func (c ooklaSpeedTestClient) Upload(url string, payload []byte) error {
	reader := bytes.NewReader(payload)

	resp, err := c.httpClient.Post(url, "application/octet-stream", reader)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))

	return nil
}
