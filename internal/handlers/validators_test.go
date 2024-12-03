package handlers

import (
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLValidate(t *testing.T) {
	type urlSlice []string

	tests := []struct {
		name     string
		urlSlice urlSlice
		want     service.Saver
		wantErr  error
	}{
		{
			name:     "Good gauge metric",
			urlSlice: urlSlice{"update", "gauge", "metricName", "1.1"},
			want:     service.GaugeMetric{Name: "metricName", Value: float64(1.1)},
			wantErr:  nil,
		},
		{
			name:     "Good counter metric",
			urlSlice: urlSlice{"update", "counter", "metricName", "1"},
			want:     service.CounterMetric{Name: "metricName", Value: int64(1)},
			wantErr:  nil,
		},
		{
			name:     "Bad gauge metric",
			urlSlice: urlSlice{"update", "gauge", "metricName", "1a"},
			want:     nil,
			wantErr:  fmt.Errorf("BadRequest"),
		},
		{
			name:     "Bad counter metric",
			urlSlice: urlSlice{"update", "counter", "metricName", "3.1"},
			want:     nil,
			wantErr:  fmt.Errorf("BadRequest"),
		},
		{
			name:     "Bad metric type",
			urlSlice: urlSlice{"update", "somethingWrong", "metricName", "3.1"},
			want:     nil,
			wantErr:  fmt.Errorf("BadRequest"),
		},
		{
			name:     "Not update in path",
			urlSlice: urlSlice{"updater", "gauge", "metricName", "3.1"},
			want:     nil,
			wantErr:  fmt.Errorf("BadRequest"),
		},
		{
			name:     "Short path",
			urlSlice: urlSlice{"updater", "gauge"},
			want:     nil,
			wantErr:  fmt.Errorf("NotFound"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLValidate(tt.urlSlice)
			assert.Equalf(t, tt.want, got, "Saver want - %v, got - %v", tt.want, got)
			assert.Equalf(t, tt.wantErr, err, "Err want - %v, got - %v", tt.wantErr, err)
		})
	}
}
