package main

import (
	"bytes"
	"github.com/goforbroke1006/speedtest-lib/domain"
	"net/http/httptest"
	"testing"
)

func Test_healthHandlerMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		wantStatusCode int
		wantContent    string
	}{
		{
			name:           "always healthy",
			wantStatusCode: 200,
			wantContent:    "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &httptest.ResponseRecorder{Body: new(bytes.Buffer)}
			healthHandlerMiddleware()(writer, nil)

			if writer.Code != tt.wantStatusCode {
				t.Errorf("healthHandlerMiddleware() code = %v, want %v", writer.Code, tt.wantStatusCode)
			}
			if writer.Body.String() != tt.wantContent {
				t.Errorf("healthHandlerMiddleware() body = %v, want %v", writer.Body.String(), tt.wantContent)
			}
		})
	}
}

func Test_readyHandlerMiddleware(t *testing.T) {
	type args struct {
		sources map[string]domain.Upgrader
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantContent    string
	}{
		{
			name: "negative 1 - all provider are not ready",
			args: args{sources: map[string]domain.Upgrader{
				"fake-1": fakeUpgrader{dl: 0, ul: 0},
				"fake-2": fakeUpgrader{dl: 0, ul: 0},
				"fake-3": fakeUpgrader{dl: 0, ul: 0},
			}},
			wantStatusCode: 404,
			wantContent:    "fail",
		},
		{
			name: "negative 2 - only one provider is not ready",
			args: args{sources: map[string]domain.Upgrader{
				"fake-1": fakeUpgrader{dl: 1, ul: 1},
				"fake-2": fakeUpgrader{dl: 3, ul: 2},
				"fake-3": fakeUpgrader{dl: 0, ul: 0},
			}},
			wantStatusCode: 404,
			wantContent:    "fail",
		},
		{
			name: "positive 1 - all providers are ready",
			args: args{sources: map[string]domain.Upgrader{
				"fake-1": fakeUpgrader{dl: 1, ul: 1},
				"fake-2": fakeUpgrader{dl: 3, ul: 2},
				"fake-3": fakeUpgrader{dl: 4, ul: 5},
			}},
			wantStatusCode: 200,
			wantContent:    "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &httptest.ResponseRecorder{Body: new(bytes.Buffer)}
			readyHandlerMiddleware(tt.args.sources)(writer, nil)

			if writer.Code != tt.wantStatusCode {
				t.Errorf("healthHandlerMiddleware() code = %v, want %v", writer.Code, tt.wantStatusCode)
			}
			if writer.Body.String() != tt.wantContent {
				t.Errorf("healthHandlerMiddleware() body = %v, want %v", writer.Body.String(), tt.wantContent)
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
