package ookla

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestNewOoklaSpeedTestClient(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "always return client instance",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOoklaSpeedTestClient(); tt.want != (got != nil) {
				t.Errorf("NewOoklaSpeedTestClient() not nil, want %v", tt.want)
			}
		})
	}
}

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
			name:    "negative - on http error",
			fields:  fields{httpClient: fakeHttpClientInt{getErr: true}},
			wantCc:  ClientConfig{},
			wantErr: true,
		},
		{
			name:    "negative - on broken resp body reader",
			fields:  fields{httpClient: fakeHttpClientInt{getBrokeRespBodyReader: true}},
			wantCc:  ClientConfig{},
			wantErr: true,
		},
		{
			name:    "negative - on broken JSON content",
			fields:  fields{httpClient: fakeHttpClientInt{getBrokeRespBodyFormat: true}},
			wantCc:  ClientConfig{},
			wantErr: true,
		},
		{
			name:   "positive - can read client config",
			fields: fields{httpClient: fakeHttpClientInt{}},
			wantCc: ClientConfig{
				XMLName: xml.Name{Local: "settings"},
				Client: Client{
					XMLName: xml.Name{Local: "client"},
					IP:      "178.70.74.118",
					Lat:     59.8983,
					Lon:     30.2618,
				},
				LicenseKey: "f7a45ced624d3a70-1df5b7cd427370f7-b91ee21d6cb22d7b",
				Download: DownloadPlan{
					XMLName:       xml.Name{Local: "download"},
					ThreadsPerUrl: 4,
				},
				Upload: UploadPlan{
					XMLName:       xml.Name{Local: "upload"},
					Threads:       2,
					ThreadsPerUrl: 4,
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
			name:     "negative - on http error",
			fields:   fields{httpClient: fakeHttpClientInt{getErr: true}},
			wantList: nil,
			wantErr:  true,
		},
		{
			name:     "negative - on broken resp body reader",
			fields:   fields{httpClient: fakeHttpClientInt{getBrokeRespBodyReader: true}},
			wantList: nil,
			wantErr:  true,
		},
		{
			name:     "negative - on broken JSON content",
			fields:   fields{httpClient: fakeHttpClientInt{getBrokeRespBodyFormat: true}},
			wantList: nil,
			wantErr:  true,
		},
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

func Test_ooklaSpeedTestClient_Download(t *testing.T) {
	type fields struct {
		httpClient httpClientInt
	}
	type args struct {
		bytesSize uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "negative - on http error",
			fields:  fields{httpClient: fakeHttpClientInt{getErr: true}},
			args:    args{bytesSize: 0},
			wantErr: true,
		},
		{
			name:    "negative - on broken resp body reader",
			fields:  fields{httpClient: fakeHttpClientInt{getBrokeRespBodyReader: true}},
			args:    args{bytesSize: 0},
			wantErr: true,
		},
		{
			name:    "negative - on wrong content length",
			fields:  fields{httpClient: fakeHttpClientInt{getCertainContentLen: 9}},
			args:    args{bytesSize: 10},
			wantErr: true,
		},
		{
			name:    "positive - correct content length",
			fields:  fields{httpClient: fakeHttpClientInt{getCertainContentLen: 111}},
			args:    args{bytesSize: 111},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ooklaSpeedTestClient{
				httpClient: tt.fields.httpClient,
			}
			if err := c.Download("tt.args.host", "tt.args.licence", tt.args.bytesSize); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ooklaSpeedTestClient_Upload(t *testing.T) {
	type fields struct {
		httpClient httpClientInt
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "negative - on http error",
			fields:  fields{httpClient: fakeHttpClientInt{postErr: true}},
			wantErr: true,
		},
		{
			name:    "negative - on broken resp body reader",
			fields:  fields{httpClient: fakeHttpClientInt{postBrokeRespBodyReader: true}},
			wantErr: true,
		},
		{
			name:    "positive - ok",
			fields:  fields{httpClient: fakeHttpClientInt{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ooklaSpeedTestClient{
				httpClient: tt.fields.httpClient,
			}
			if err := c.Upload("http://some-server/upload.php", []byte("tt.args.payload")); (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var (
	_ httpClientInt = &fakeHttpClientInt{}
)

type fakeHttpClientInt struct {
	getErr                  bool
	postErr                 bool
	getBrokeRespBodyReader  bool
	postBrokeRespBodyReader bool
	getBrokeRespBodyFormat  bool
	getCertainContentLen    uint
}

func (f fakeHttpClientInt) Get(url string) (resp *http.Response, err error) {
	if f.getErr {
		return nil, errors.New("fake error")
	}

	if f.getBrokeRespBodyReader {
		resp = new(http.Response)
		resp.Body = io.NopCloser(brokenReader{})
		return
	}

	if f.getCertainContentLen > 0 {
		resp = new(http.Response)
		buf := make([]byte, f.getCertainContentLen)
		resp.Body = io.NopCloser(strings.NewReader(string(buf)))
		return
	}

	if strings.Contains(url, "speedtest-config.php") {
		resp = new(http.Response)
		readFile, err := ioutil.ReadFile("./testdata/speedtest-config-1.xml")
		if err != nil {
			panic(err)
		}
		content := string(readFile)

		if f.getBrokeRespBodyFormat {
			content = "wildfowl for broken content" + content
			content = content[:len(content)-1]
		}

		resp.Body = io.NopCloser(strings.NewReader(content))

		return resp, nil
	}

	if strings.Contains(url, "api/js/servers") {
		resp = new(http.Response)
		readFile, err := ioutil.ReadFile("./testdata/servers-1.json")
		if err != nil {
			panic(err)
		}
		content := string(readFile)

		if f.getBrokeRespBodyFormat {
			content = "wildfowl for broken content" + content
			content = content[:len(content)-1]
		}

		resp.Body = io.NopCloser(strings.NewReader(content))

		return resp, nil
	}

	panic("implement me")
}

func (f fakeHttpClientInt) Post(url, _ string, _ io.Reader) (resp *http.Response, err error) {
	if f.postErr {
		return nil, errors.New("fake error")
	}

	if f.postBrokeRespBodyReader {
		resp = new(http.Response)
		resp.Body = io.NopCloser(brokenReader{})
		return
	}

	if strings.Contains(url, "upload.php") {
		resp = new(http.Response)
		resp.Body = io.NopCloser(strings.NewReader("hello"))
		return
	}

	panic("implement me")
}

type brokenReader struct {
}

func (br brokenReader) Read(_ []byte) (int, error) {
	return 0, errors.New("fake error")
}
