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
		name    string
		mType   string
		mName   string
		want    service.MetricHolder
		wantErr error
	}{
		{
			name:    "Good gauge metric",
			mType:   "gauge",
			mName:   "metricName",
			want:    &service.GaugeMetric{Name: "metricName", Value: float64(0)},
			wantErr: nil,
		},
		{
			name:    "Good counter metric",
			mType:   "counter",
			mName:   "metricName",
			want:    &service.CounterMetric{Name: "metricName", Value: int64(0)},
			wantErr: nil,
		},
		{
			name:    "Bad metric gauge type",
			mType:   "gaugeR",
			mName:   "metricName",
			want:    nil,
			wantErr: fmt.Errorf("BadRequest"),
		},
		{
			name:    "Bad metric counter type",
			mType:   "counterR",
			mName:   "metricName",
			want:    nil,
			wantErr: fmt.Errorf("BadRequest"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLValidate(tt.mType, tt.mName)

			assert.Equalf(t, tt.want, got, "Saver want - %v, got - %v", tt.want, got)
			assert.Equalf(t, tt.wantErr, err, "Err want - %v, got - %v", tt.wantErr, err)
		})
	}
}
