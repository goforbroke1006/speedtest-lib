package content_len

import (
	"reflect"
	"testing"
)

func TestDataLen(t *testing.T) {
	type args struct {
		bytesCount float64
	}
	tests := []struct {
		name string
		args args
		want *converter
	}{
		{
			name: "sample 1",
			args: args{bytesCount: 0},
			want: &converter{bytesCount: 0},
		},
		{
			name: "sample 2",
			args: args{bytesCount: 1024},
			want: &converter{bytesCount: 1024},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DataLen(tt.args.bytesCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_converter_Bites(t *testing.T) {
	type fields struct {
		bytesCount float64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name:   "0 byte",
			fields: fields{bytesCount: 0},
			want:   0,
		},
		{
			name:   "1 byte",
			fields: fields{bytesCount: 1},
			want:   8,
		},
		{
			name:   "half of byte",
			fields: fields{bytesCount: 0.5},
			want:   4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := converter{
				bytesCount: tt.fields.bytesCount,
			}
			if got := c.Bites(); got != tt.want {
				t.Errorf("Bites() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_converter_MegaBites validates to-megabits conversion
// Actual results are checked with https://convertlive.com/u/convert/bytes/to/megabits
func Test_converter_MegaBites(t *testing.T) {
	type fields struct {
		bytesCount float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "1 kilobyte",
			fields: fields{bytesCount: 1024},
			want:   0.0078125,
		},
		{
			name:   "4 kilobyte",
			fields: fields{bytesCount: 4096},
			want:   0.03125,
		},
		{
			name:   "1 megabyte",
			fields: fields{bytesCount: 1048576},
			want:   8.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := converter{
				bytesCount: tt.fields.bytesCount,
			}
			if got := c.MegaBites(); got != tt.want {
				t.Errorf("MegaBites() = %v, want %v", got, tt.want)
			}
		})
	}
}
