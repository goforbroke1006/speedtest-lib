package speedtest_lib

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"unsafe"

	"github.com/golang/mock/gomock"
	"github.com/showwin/speedtest-go/speedtest"

	"github.com/goforbroke1006/speedtest-lib/domain"
	"github.com/goforbroke1006/speedtest-lib/mocks"
)

func Test_newOoklaLoader(t *testing.T) {
	got := newOoklaLoader()
	if got == nil {
		t.Errorf("newOoklaLoader() = nil")
	}
}

func Test_ooklaLoaderDefault_LoadConfig(t *testing.T) {
	type fields struct {
		client  func(ctrl *gomock.Controller) domain.OoklaSpeedTestClient
		targets speedtest.Servers
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "negative 1 - FetchUserInfo failed",
			fields: fields{
				client: func(ctrl *gomock.Controller) domain.OoklaSpeedTestClient {
					m := mocks.NewMockOoklaSpeedTestClient(ctrl)
					m.EXPECT().FetchUserInfo().Return(nil, errors.New("fake error"))
					return m
				},
			},
			wantErr: true,
		},
		{
			name: "negative 1 - FetchServers failed",
			fields: fields{
				client: func(ctrl *gomock.Controller) domain.OoklaSpeedTestClient {
					m := mocks.NewMockOoklaSpeedTestClient(ctrl)
					m.EXPECT().FetchUserInfo().Return(new(speedtest.User), nil)
					m.EXPECT().FetchServers(gomock.Any()).Return(nil, errors.New("fake error"))
					return m
				},
			},
			wantErr: true,
		},
		{
			name: "negative 1 - FetchServers works fine but returns empty list",
			fields: fields{
				client: func(ctrl *gomock.Controller) domain.OoklaSpeedTestClient {
					m := mocks.NewMockOoklaSpeedTestClient(ctrl)
					m.EXPECT().FetchUserInfo().Return(new(speedtest.User), nil)
					m.EXPECT().FetchServers(gomock.Any()).Return(speedtest.Servers{}, nil)
					return m
				},
			},
			wantErr: true,
		},
		{
			name: "positive",
			fields: fields{
				client: func(ctrl *gomock.Controller) domain.OoklaSpeedTestClient {
					m := mocks.NewMockOoklaSpeedTestClient(ctrl)
					m.EXPECT().FetchUserInfo().Return(new(speedtest.User), nil)
					m.EXPECT().FetchServers(gomock.Any()).Return(speedtest.Servers{
						&speedtest.Server{},
						&speedtest.Server{},
					}, nil)
					return m
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			client := tt.fields.client(mockCtrl)

			o := &ooklaLoaderDefault{
				client:  client,
				targets: tt.fields.targets,
			}
			if err := o.LoadConfig(); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ooklaLoaderDefault_DownloadSink(t *testing.T) {
	type fields struct {
		client  domain.OoklaSpeedTestClient
		targets speedtest.Servers
	}
	tests := []struct {
		name    string
		fields  fields
		want    []float64
		wantErr bool
	}{
		{
			name: "positive - no servers = no data, no errors",
			fields: fields{
				targets: speedtest.Servers{},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := ooklaLoaderDefault{
				client:  tt.fields.client,
				targets: tt.fields.targets,
			}

			var got []float64
			sink, err := o.DownloadSink()

			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for c := range sink {
				got = append(got, c)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DownloadSink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: have to use mock for speedtest.Server to prevent network usage
func Test_ooklaLoaderDefault_DownloadSink_callWithNetwork(t *testing.T) {
	fixServerDoer := func(s *speedtest.Server) {
		// need to set private field to allow Pint operation invocation
		serverDoer := reflect.ValueOf(s).Elem().FieldByName("doer")
		ptrToDoer := unsafe.Pointer(serverDoer.UnsafeAddr())
		realPtrToY := (**http.Client)(ptrToDoer)
		*realPtrToY = http.DefaultClient
	}

	serverOne := &speedtest.Server{URL: "http://sankt-peterburg2.speedtest.rt.ru:8080/speedtest/upload.php"}
	fixServerDoer(serverOne)

	serversList := speedtest.Servers{
		serverOne,
	}

	o := ooklaLoaderDefault{
		targets: serversList,
	}

	var got []float64
	sink, err := o.DownloadSink()

	if err != nil {
		t.Errorf("DownloadSink() error should be NIL")
		return
	}

	for c := range sink {
		got = append(got, c)
	}

	if len(got) != 1 {
		t.Errorf("DownloadSink() got 1 answer from 1 server, got = %v", got)
	}
}

func Test_ooklaLoaderDefault_UploadSink(t *testing.T) {
	type fields struct {
		client  domain.OoklaSpeedTestClient
		targets speedtest.Servers
	}
	tests := []struct {
		name    string
		fields  fields
		want    []float64
		wantErr bool
	}{
		{
			name: "positive - no servers = no data, no errors",
			fields: fields{
				targets: speedtest.Servers{},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := ooklaLoaderDefault{
				client:  tt.fields.client,
				targets: tt.fields.targets,
			}

			var got []float64
			sink, err := o.UploadSink()

			if (err != nil) != tt.wantErr {
				t.Errorf("UploadSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for c := range sink {
				got = append(got, c)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploadSink() got = %v, want %v", got, tt.want)
			}
		})
	}
}
