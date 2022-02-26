package test_speed

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_requestValidator_Validate(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name     string
		args     args
		wantErrs []string
	}{
		{
			name: "negative 1 - no provider param",
			args: args{req: func() *http.Request {
				request, err := http.NewRequest(http.MethodGet, "https://some-host.com/test?some=param", nil)
				if err != nil {
					t.Fatal(err)
				}
				return request
			}()},
			wantErrs: []string{
				"'provider' param is required",
			},
		},
		{
			name: "negative 2 - has provider param but empty",
			args: args{req: func() *http.Request {
				request, err := http.NewRequest(http.MethodGet, "https://some-host.com/test?provider=", nil)
				if err != nil {
					t.Fatal(err)
				}
				return request
			}()},
			wantErrs: []string{
				"'provider' param should not be empty",
			},
		},
		{
			name: "positive 1 - has provider",
			args: args{req: func() *http.Request {
				request, err := http.NewRequest(http.MethodGet, "https://some-host.com/test?provider=any", nil)
				if err != nil {
					t.Fatal(err)
				}
				return request
			}()},
			wantErrs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := requestValidator{}
			if gotErrs := rv.Validate(tt.args.req); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("Validate() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}
