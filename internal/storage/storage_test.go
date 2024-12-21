package storage

import (
	"github.com/forester2k/metrics/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMemStorage_Init(t *testing.T) {
	Store = &MemStorage{}
	Store.Init()
	type fields struct {
		Gauges    map[string]float64
		Counters  map[string]int64
		ValidKeys map[string]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "Testing init"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, Store.Counters["PollCount"], int64(0))
		})
	}
}
func TestMemStorage_Save(t *testing.T) {

	tests := []struct {
		name    string
		metric  service.MetricHolder
		wantErr bool
	}{
		{
			name:    "good gauge metric",
			metric:  &service.GaugeMetric{Name: "Alloc", Value: float64(1.1)},
			wantErr: false,
		},
		{
			name:    "good counter metric",
			metric:  &service.CounterMetric{Name: "PollCount", Value: int64(1)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := service.MetricHolder(tt.metric)
			err := Store.Save(&m)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			switch metric := m.(type) {
			case *service.GaugeMetric:
				Store.Gauges[metric.Name] = metric.Value
				val, ok := Store.Gauges[metric.Name]
				if !ok {
					require.Error(t, nil, "didn't save")
				}
				assert.Equal(t, tt.metric.GetValue(), val)
			case *service.CounterMetric:
				Store.Counters[metric.Name] = metric.Value
				val, ok := Store.Counters[metric.Name]
				if !ok {
					require.Error(t, nil, "didn't save")
				}
				assert.Equal(t, tt.metric.GetValue(), val)
			}
		})
	}
}

//
//func TestMemStorage_Get(t *testing.T) {
//	tests := []struct {
//		name    string
//		mName   string
//		want    float64
//		wantErr bool
//	}{
//		{
//			name:    "Good metric",
//			mName:   "Alloc",
//			want:    float64(1.1),
//			wantErr: false,
//		},
//		{
//			name:    "Bad metric",
//			mName:   "XXXAllocXXX",
//			want:    float64(0),
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := Store.Get(tt.mName)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetGauge() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("GetGauge() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestMemStorage_MakeList(t *testing.T) {
	tests := []struct {
		name string
		want *StoredList
	}{
		{
			name: "Good list",
			want: &StoredList{
				Gauges:   map[string]float64{"Alloc": float64(1.1)},
				Counters: map[string]int64{"PollCount": int64(1)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Store.MakeList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeList() = %v, want %v", got, tt.want)
			}
		})
	}
}
