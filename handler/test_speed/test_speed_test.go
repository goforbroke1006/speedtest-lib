package test_speed

import (
	"github.com/goforbroke1006/speedtest-lib/upgrader"
	"testing"
)

func TestTestSpeedHandler(t *testing.T) {
	type args struct {
		sources map[string]upgrader.Upgrader
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
			if got := TestSpeedHandler(tt.args.sources); tt.want != (got != nil) {
				t.Errorf("TestSpeedHandler() want %v", tt.want)
			}
		})
	}
}
