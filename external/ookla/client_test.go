package ookla

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_ooklaSpeedTestClient_GetClientConfig(t *testing.T) {
	type fields struct {
		httpClient httpClientInt
	}
	tests := []struct {
		name    string
		fields  fields
		wantCc  ClientConfig
		wantErr bool
	}{
		{
			name:   "can read client config",
			fields: fields{httpClient: fakeHttpClientInt{}},
			wantCc: ClientConfig{
				XMLName: xml.Name{Local: "settings"},
				Client: Client{
					XMLName: xml.Name{Local: "client"},
					IP:      "178.70.74.118",
					Lat:     "59.8983",
					Lon:     "30.2618",
				},
				LicenseKey: "f7a45ced624d3a70-1df5b7cd427370f7-b91ee21d6cb22d7b",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ooklaSpeedTestClient{
				httpClient: tt.fields.httpClient,
			}
			gotCc, err := c.GetClientConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClientConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCc, tt.wantCc) {
				t.Errorf("GetClientConfig() gotCc = %v, want %v", gotCc, tt.wantCc)
			}
		})
	}
}

func Test_ooklaSpeedTestClient_GetServersList(t *testing.T) {
	type fields struct {
		httpClient httpClientInt
	}
	tests := []struct {
		name     string
		fields   fields
		wantList ServersList
		wantErr  bool
	}{
		{
			name:   "positive 1 - with fake response",
			fields: fields{httpClient: fakeHttpClientInt{}},
			wantList: ServersList{
				ServerSummary{
					ID:       "2599",
					Name:     "Saint Petersburg",
					Url:      "http://sankt-peterburg2.speedtest.rt.ru:8080/speedtest/upload.php",
					Host:     "sankt-peterburg2.speedtest.rt.ru.prod.hosts.ooklaserver.net:8080",
					Lat:      "59.9333",
					Lon:      "30.3333",
					Distance: 5,
					Country:  "Russia",
					CC:       "RU",
					Sponsor:  "Rostelecom",
				},
				ServerSummary{
					ID:       "28732",
					Name:     "Saint Petersburg",
					Url:      "http://speedtest.viltel.net:8080/speedtest/upload.php",
					Host:     "speedtest.viltel.net:8080",
					Lat:      "59.9343",
					Lon:      "30.3351",
					Distance: 5,
					Country:  "Russia",
					CC:       "RU",
					Sponsor:  "Viltel Ltd",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ooklaSpeedTestClient{
				httpClient: tt.fields.httpClient,
			}
			gotList, err := c.GetServersList()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServersList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotList, tt.wantList) {
				t.Errorf("GetServersList() gotList = %v, want %v", gotList, tt.wantList)
			}
		})
	}
}

var (
	_ httpClientInt = &fakeHttpClientInt{}
)

type fakeHttpClientInt struct {
}

func (f fakeHttpClientInt) Get(url string) (resp *http.Response, err error) {
	if strings.Contains(url, "speedtest-config.php") {
		resp = new(http.Response)
		readFile, err := ioutil.ReadFile("./testdata/speedtest-config-1.xml")
		if err != nil {
			panic(err)
		}
		resp.Body = io.NopCloser(strings.NewReader(string(readFile)))
		return resp, nil
	}

	if strings.Contains(url, "api/js/servers") {
		resp = new(http.Response)
		readFile, err := ioutil.ReadFile("./testdata/servers-1.json")
		if err != nil {
			panic(err)
		}
		resp.Body = io.NopCloser(strings.NewReader(string(readFile)))
		return resp, nil
	}

	panic("implement me")
}

func (f fakeHttpClientInt) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	//TODO implement me
	panic("implement me")
}
