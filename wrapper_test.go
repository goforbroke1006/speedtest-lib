package speedtest_lib

import (
	"context"
	"testing"
)

func TestDoRequest(t *testing.T) {
	type args struct {
		ctx  context.Context
		kind ProviderKind
	}
	tests := []struct {
		name         string
		args         args
		wantDownload float64
		wantUpload   float64
		wantErr      bool
	}{
		{
			name: "negative 1 - unknown provider",
			args: args{
				ctx:  context.Background(),
				kind: "wildfowl",
			},
			wantDownload: 0,
			wantUpload:   0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapper{}
			gotDownload, gotUpload, err := w.DoRequest(tt.args.ctx, tt.args.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDownload != tt.wantDownload {
				t.Errorf("DoRequest() gotDownload = %v, want %v", gotDownload, tt.wantDownload)
			}
			if gotUpload != tt.wantUpload {
				t.Errorf("DoRequest() gotUpload = %v, want %v", gotUpload, tt.wantUpload)
			}
		})
	}
}

// Benchmark_DoRequest_Ookla-12    	       1	15839082397 ns/op
func Benchmark_DoRequest_Ookla(b *testing.B) {
	w := New()
	for i := 0; i < b.N; i++ {
		download, upload, err := w.DoRequest(context.Background(), ProviderKindOokla)
		if err != nil {
			b.Fatal(err)
		}
		_, _ = download, upload
	}
}

// Benchmark_DoRequest_Netflix-12    	       1	12117915392 ns/op
func Benchmark_DoRequest_Netflix(b *testing.B) {
	w := New()
	for i := 0; i < b.N; i++ {
		download, upload, err := w.DoRequest(context.Background(), ProviderKindNetflix)
		if err != nil {
			b.Fatal(err)
		}
		_, _ = download, upload
	}
}
