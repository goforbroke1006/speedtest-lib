package content

import (
	"math"
	"testing"
)

func TestDataLen_Bits(t *testing.T) {
	tests := []struct {
		name string
		dl   DataLen
		want uint64
	}{
		{
			name: "positive 1 - zero",
			dl:   DataLen(0),
			want: 0,
		},
		{
			name: "positive 2 - 1 bit",
			dl:   DataLen(1 * Bit),
			want: 1,
		},
		{
			name: "positive 3 - 1 byte",
			dl:   DataLen(1 * Byte),
			want: 8,
		},
		{
			name: "positive 4 - 1 mega bit",
			dl:   DataLen(1 * MegaBit),
			want: 1048576,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dl.Bits(); got != tt.want {
				t.Errorf("Bits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataLen_Bytes(t *testing.T) {
	tests := []struct {
		name string
		dl   DataLen
		want float64
	}{
		{
			name: "positive - zero",
			dl:   DataLen(0),
			want: 0,
		},
		{
			name: "positive - 1 bit",
			dl:   DataLen(1 * Bit),
			want: 0.125,
		},
		{
			name: "positive - 1 byte",
			dl:   DataLen(1 * Byte),
			want: 1.0,
		},
		{
			name: "positive - 1 mega bit",
			dl:   DataLen(1 * MegaBit),
			want: 131072,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dl.Bytes(); got != tt.want {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataLen_KiloBytes(t *testing.T) {
	tests := []struct {
		name string
		dl   DataLen
		want float64
	}{
		{
			name: "zero",
			dl:   DataLen(0),
			want: 0,
		},
		{
			name: "1024 bytes",
			dl:   DataLen(1024 * Byte),
			want: 1,
		},
		{
			name: "1 kilo bytes",
			dl:   DataLen(1 * KiloByte),
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dl.KiloBytes(); got != tt.want {
				t.Errorf("KiloBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataLen_MegaBites(t *testing.T) {
	tests := []struct {
		name string
		dl   DataLen
		want float64
	}{
		{
			name: "zero",
			dl:   DataLen(0),
			want: 0,
		},
		{
			name: "byte",
			dl:   DataLen(1 * Byte),
			want: 0.0000076294,
		},
		{
			name: "1 mega bit",
			dl:   DataLen(1 * MegaBit),
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dl.MegaBites(); math.Abs(got-tt.want) > 0.0001 {
				t.Errorf("MegaBites() = %v, want %v", got, tt.want)
			}
		})
	}
}
