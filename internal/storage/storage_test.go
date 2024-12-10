package storage

import (
	"github.com/stretchr/testify/assert"
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

//func Test_isValidMetric(t *testing.T) {
//	type args struct {
//		name  string
//		mType string
//		m     *MemStorage
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "valid metric",
//			args: struct {
//				name  string
//				mType string
//				m     *MemStorage
//			}{
//				"PollCount",
//				"counter",
//				Store,
//			},
//			wantErr: false,
//		},
//		{
//			name: "invalid metric",
//			args: struct {
//				name  string
//				mType string
//				m     *MemStorage
//			}{
//				"XXXPollCountXXX",
//				"counter",
//				Store,
//			},
//			wantErr: true,
//		},
//		{
//			name: "invalid metric type",
//			args: struct {
//				name  string
//				mType string
//				m     *MemStorage
//			}{
//				"PollCount",
//				"gauge",
//				Store,
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := isValidMetric(tt.args.name, tt.args.mType, tt.args.m); (err != nil) != tt.wantErr {
//				t.Errorf("isValidMetric() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestMemStorage_AddGauge(t *testing.T) {
	tests := []struct {
		name    string
		mName   string
		mValue  float64
		wantErr bool
	}{
		{
			name:    "good metric",
			mName:   "Alloc",
			mValue:  float64(1.1),
			wantErr: false,
		},
		//{
		//	name:    "bad metric",
		//	mName:   "XXXAllocXXX",
		//	mValue:  float64(2.2),
		//	wantErr: true,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Store.AddGauge(tt.mName, tt.mValue); (err != nil) != tt.wantErr {
				t.Errorf("AddGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	tests := []struct {
		name    string
		mName   string
		mValue  int64
		wantErr bool
	}{
		{
			name:    "good metric",
			mName:   "PollCount",
			mValue:  int64(1),
			wantErr: false,
		},
		//{
		//	name:    "bad metric",
		//	mName:   "XXXPollCountXXX",
		//	mValue:  int64(2),
		//	wantErr: true,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Store.AddCounter(tt.mName, tt.mValue); (err != nil) != tt.wantErr {
				t.Errorf("AddGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	tests := []struct {
		name    string
		mName   string
		want    float64
		wantErr bool
	}{
		{
			name:    "Good metric",
			mName:   "Alloc",
			want:    float64(1.1),
			wantErr: false,
		},
		{
			name:    "Bad metric",
			mName:   "XXXAllocXXX",
			want:    float64(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Store.GetGauge(tt.mName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGauge() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	tests := []struct {
		name    string
		mName   string
		want    int64
		wantErr bool
	}{
		{
			name:    "Good metric",
			mName:   "PollCount",
			want:    int64(1),
			wantErr: false,
		},
		{
			name:    "Bad metric",
			mName:   "XXXPollCountXXX",
			want:    int64(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Store.GetCounter(tt.mName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGauge() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Store.MakeList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeList() = %v, want %v", got, tt.want)
			}
		})
	}
}
