package speedtest_lib

import (
	"context"
	"github.com/goforbroke1006/speedtest-lib/domain"
	"testing"
)

func Test_New(t *testing.T) {
	instance := New()
	if instance == nil {
		t.Errorf("New() should not be NIL")
	}
}

func Test_wrapper_DoRequest(t *testing.T) {
	type fields struct {
		loaderOokla   domain.NetworkLoader
		loaderNetflix domain.NetworkLoader
	}
	type args struct {
		ctx  context.Context
		kind ProviderKind
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantDownload float64
		wantUpload   float64
		wantErr      bool
	}{
		{
			name:   "negative 1 - unknown provider",
			fields: fields{},
			args: args{
				ctx:  context.Background(),
				kind: "wildfowl",
			},
			wantDownload: 0,
			wantUpload:   0,
			wantErr:      true,
		},
		{
			name: "positive 1 - fake ookla",
			fields: fields{
				loaderOokla: fakeLoader{},
			},
			args: args{
				ctx:  context.Background(),
				kind: "ookla",
			},
			wantDownload: 1,
			wantUpload:   1,
			wantErr:      false,
		},
		{
			name: "positive 1 - fake netflix",
			fields: fields{
				loaderNetflix: fakeLoader{},
			},
			args: args{
				ctx:  context.Background(),
				kind: "netflix",
			},
			wantDownload: 1,
			wantUpload:   1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wrapper{
				loaderOokla:   tt.fields.loaderOokla,
				loaderNetflix: tt.fields.loaderNetflix,
			}
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

type fakeLoader struct{}

var _ domain.NetworkLoader = &fakeLoader{}

func (l fakeLoader) LoadConfig() error {
	return nil
}

func (l fakeLoader) DownloadSink() (bits <-chan float64, err error) {
	bitsPerSecondSink := make(chan float64)
	go func() {
		bitsPerSecondSink <- 1
		close(bitsPerSecondSink)
	}()
	return bitsPerSecondSink, nil
}

func (l fakeLoader) UploadSink() (bits <-chan float64, err error) {
	bitsPerSecondSink := make(chan float64)
	go func() {
		bitsPerSecondSink <- 1
		close(bitsPerSecondSink)
	}()
	return bitsPerSecondSink, nil
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
