package speedtest_lib

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/goforbroke1006/speedtest-lib/domain"
	"github.com/goforbroke1006/speedtest-lib/mocks"
)

func Test_newNetflixLoader(t *testing.T) {
	got := newNetflixLoader()
	if got == nil {
		t.Errorf("newNetflixLoader() = nil")
	}
}

func Test_netflixLoader_LoadConfig(t *testing.T) {
	type fields struct {
		fastCom func(ctrl *gomock.Controller) domain.NetflixFastClient
		urls    []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "negative",
			fields: fields{
				fastCom: func(ctrl *gomock.Controller) domain.NetflixFastClient {
					m := mocks.NewMockNetflixFastClient(ctrl)
					m.EXPECT().GetUrls().Return(nil, errors.New("fake error"))
					return m
				},
				urls: nil,
			},
			wantErr: true,
		},
		{
			name: "positive",
			fields: fields{
				fastCom: func(ctrl *gomock.Controller) domain.NetflixFastClient {
					m := mocks.NewMockNetflixFastClient(ctrl)
					m.EXPECT().GetUrls().Return([]string{}, nil)
					return m
				},
				urls: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fastCom := tt.fields.fastCom(mockCtrl)

			n := &netflixLoader{
				fastCom: fastCom,
				urls:    tt.fields.urls,
			}
			if err := n.LoadConfig(); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_netflixLoader_DownloadSink(t *testing.T) {
	type fields struct {
		fastCom func(ctrl *gomock.Controller) domain.NetflixFastClient
		urls    []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []float64
		wantErr bool
	}{
		{
			name: "positive - no data",
			fields: fields{
				fastCom: func(ctrl *gomock.Controller) domain.NetflixFastClient {
					var (
						urls []string
						ch   = make(chan<- float64)
					)

					m := mocks.NewMockNetflixFastClient(ctrl)
					m.EXPECT().
						Measure(gomock.AssignableToTypeOf(urls), gomock.AssignableToTypeOf(ch)).
						DoAndReturn(
							func(urls []string, KbpsChan chan<- float64) error {
								go func() {
									// no data sent
									close(KbpsChan)
								}()
								return nil
							},
						)
					return m
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "positive - regular",
			fields: fields{
				fastCom: func(ctrl *gomock.Controller) domain.NetflixFastClient {
					var (
						urls []string
						ch   = make(chan<- float64)
					)

					m := mocks.NewMockNetflixFastClient(ctrl)
					m.EXPECT().
						Measure(gomock.AssignableToTypeOf(urls), gomock.AssignableToTypeOf(ch)).
						DoAndReturn(
							func(urls []string, KbpsChan chan<- float64) error {
								go func() {
									KbpsChan <- 1.0 // 1 Kb = 1024 bit
									KbpsChan <- 2.0 // 2 Kb = 2048 bit
									close(KbpsChan)
								}()
								return nil
							},
						)
					return m
				},
			},
			want:    []float64{1024.0, 2048.0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fastCom := tt.fields.fastCom(mockCtrl)

			n := netflixLoader{
				fastCom: fastCom,
				urls:    tt.fields.urls,
			}

			sink, err := n.DownloadSink()

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			var got []float64
			for c := range sink {
				got = append(got, c)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DownloadSink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_netflixLoader_UploadSink(t *testing.T) {
	type fields struct {
		fastCom domain.NetflixFastClient
		urls    []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []uint64
		wantErr bool
	}{
		{
			name:    "positive - stub works fine",
			fields:  fields{},
			want:    []uint64{0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := netflixLoader{
				fastCom: tt.fields.fastCom,
				urls:    tt.fields.urls,
			}

			sink, err := n.UploadSink()

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			var got []float64
			for c := range sink {
				got = append(got, c)
			}

			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploadSink() = %v, want %v", got, tt.want)
			}
		})
	}
}
