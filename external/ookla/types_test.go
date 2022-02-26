package ookla

import "testing"

func TestServerSummary_GetID(t *testing.T) {
	type fields struct {
		ID string
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name:   "positive 1 - simple number",
			fields: fields{ID: "138727328"},
			want:   138727328,
		},
		{
			name:   "positive 2 - zero",
			fields: fields{ID: "0"},
			want:   0,
		},
		{
			name:   "negative 1 - invalid",
			fields: fields{ID: "hello world"},
			want:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := ServerSummary{
				ID: tt.fields.ID,
			}
			if got := ss.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerSummary_GetLat(t *testing.T) {
	type fields struct {
		Lat string
	}
	var tests = []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "positive 1",
			fields: fields{Lat: "0.0"},
			want:   0,
		},
		{
			name:   "positive 2",
			fields: fields{Lat: "0.1"},
			want:   0.1,
		},
		{
			name:   "positive 3",
			fields: fields{Lat: "59.8983"},
			want:   59.8983,
		},
		{
			name:   "negative 1",
			fields: fields{Lat: "hello world"},
			want:   0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := ServerSummary{
				Lat: tt.fields.Lat,
			}
			if got := ss.GetLat(); got != tt.want {
				t.Errorf("GetLat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerSummary_GetLon(t *testing.T) {
	type fields struct {
		Lon string
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "positive 1",
			fields: fields{Lon: "0.0"},
			want:   0,
		},
		{
			name:   "positive 2",
			fields: fields{Lon: "0.1"},
			want:   0.1,
		},
		{
			name:   "positive 3",
			fields: fields{Lon: "30.2618"},
			want:   30.2618,
		},
		{
			name:   "negative 1",
			fields: fields{Lon: "hello world"},
			want:   0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := ServerSummary{
				Lon: tt.fields.Lon,
			}
			if got := ss.GetLon(); got != tt.want {
				t.Errorf("GetLon() = %v, want %v", got, tt.want)
			}
		})
	}
}
