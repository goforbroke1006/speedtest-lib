package test_speed

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/goforbroke1006/speedtest-lib/domain"
)

func Test_HandlerMiddleware(t *testing.T) {
	type args struct {
		sources map[string]domain.Upgrader
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return with any args",
			args: args{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandlerMiddleware(tt.args.sources); tt.want != (got != nil) {
				t.Errorf("HandlerMiddleware() want %v", tt.want)
			}
		})
	}
}

func Test_HandlerMiddleware_Func(t *testing.T) {
	type args struct {
		sources map[string]domain.Upgrader
		req     *http.Request
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "positive 1 - wrong param",
			args: args{
				sources: nil,
				req: func() *http.Request {
					request, err := http.NewRequest(http.MethodGet, "https://some-host.com/test?some=param", nil)
					if err != nil {
						t.Fatal(err)
					}
					return request
				}(),
			},
			wantStatusCode: 400,
		},
		{
			name: "positive 2 - empty param",
			args: args{
				sources: nil,
				req: func() *http.Request {
					request, err := http.NewRequest(http.MethodGet, "https://some-host.com/test?provider=", nil)
					if err != nil {
						t.Fatal(err)
					}
					return request
				}(),
			},
			wantStatusCode: 400,
		},
		{
			name: "positive 3 - param ok, but source not found",
			args: args{
				sources: nil,
				req: func() *http.Request {
					request, err := http.NewRequest(http.MethodGet, "https://some-host.com/test?provider=wildfowl", nil)
					if err != nil {
						t.Fatal(err)
					}
					return request
				}(),
			},
			wantStatusCode: 400,
		},
		{
			name: "positive 4 - param ok, source ok",
			args: args{
				sources: map[string]domain.Upgrader{
					"wildfowl": fakeUpgrader{},
				},
				req: func() *http.Request {
					request, err := http.NewRequest(http.MethodGet, "https://some-host.com/test?provider=wildfowl", nil)
					if err != nil {
						t.Fatal(err)
					}
					return request
				}(),
			},
			wantStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HandlerMiddleware(tt.args.sources)
			writer := &httptest.ResponseRecorder{Body: new(bytes.Buffer)}

			handler(writer, tt.args.req)

			if tt.wantStatusCode != writer.Code {
				t.Errorf("HandlerMiddleware()() statusCode, got =%v want %v", tt.wantStatusCode, writer.Code)
			}
		})
	}
}

var (
	_ domain.Upgrader = &fakeUpgrader{}
)

type fakeUpgrader struct {
	dl, ul float64
}

func (f fakeUpgrader) Run() {
}

func (f fakeUpgrader) IsReady() bool {
	return f.dl > 0 && f.ul > 0
}

func (f fakeUpgrader) GetDLSpeedMbps() float64 {
	return f.dl
}

func (f fakeUpgrader) GetULSpeedMbps() float64 {
	return f.ul
}

// BenchmarkHandlerMiddlewareFunc-12        1269525               805.8 ns/op
func BenchmarkHandlerMiddlewareFunc(b *testing.B) {
	sources := map[string]domain.Upgrader{
		"wildfowl": fakeUpgrader{},
	}
	handler := HandlerMiddleware(sources)

	request, err := http.NewRequest(http.MethodGet, "https://some-host.com/test?provider=wildfowl", nil)
	if err != nil {
		b.Fatal(err)
	}
	writer := &httptest.ResponseRecorder{Body: new(bytes.Buffer)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler(writer, request)

		if writer.Code != http.StatusOK {
			b.FailNow()
		}
	}
}

var (
	_ domain.Upgrader = &fakeUpgraderWithMutex{}
)

type fakeUpgraderWithMutex struct {
	dl, ul float64
	mx     sync.RWMutex
}

func (f *fakeUpgraderWithMutex) Run() {
}

func (f *fakeUpgraderWithMutex) IsReady() bool {
	f.mx.RLock()
	defer f.mx.RUnlock()

	return f.dl > 0 && f.ul > 0
}

func (f *fakeUpgraderWithMutex) GetDLSpeedMbps() float64 {
	f.mx.RLock()
	defer f.mx.RUnlock()

	return f.dl
}

func (f *fakeUpgraderWithMutex) GetULSpeedMbps() float64 {
	f.mx.RLock()
	defer f.mx.RUnlock()

	return f.ul
}
