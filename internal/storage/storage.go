package storage

import (
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"sync"
)

type MemStorage struct {
	GMx      sync.Mutex         `json:"-"`
	CMx      sync.Mutex         `json:"-"`
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]int64   `json:"counters"`
}

func (store *MemStorage) Save(m *service.MetricHolder) error {
	switch metric := (*m).(type) {
	case *service.GaugeMetric:
		store.GMx.Lock()
		store.Gauges[metric.Name] = metric.Value
		store.GMx.Unlock()
		StoreSynqSave <- true
		return nil
	case *service.CounterMetric:
		store.CMx.Lock()
		if _, ok := store.Counters[metric.Name]; !ok {
			store.Counters[metric.Name] = 0
		}
		store.Counters[metric.Name] = store.Counters[metric.Name] + metric.Value
		store.CMx.Unlock()
		StoreSynqSave <- true
		return nil
	default:
		return fmt.Errorf("func (store *MemStorage) Save: this should never happen")
	}
}

var Store *MemStorage
var StoreSynqSave chan bool

func (store *MemStorage) Init() {
	store.Gauges = make(map[string]float64)
	store.Counters = make(map[string]int64)
	store.Counters["PollCount"] = 0
	StoreSynqSave = make(chan bool)
}

func (store *MemStorage) GetValue(m *service.MetricHolder) (any, error) {
	switch metric := (*m).(type) {
	case *service.GaugeMetric:
		store.GMx.Lock()
		value, ok := store.Gauges[metric.Name]
		store.GMx.Unlock()
		if ok {
			//metric.Value = value
			return value, nil
		}
		return 0, fmt.Errorf("storage: metric %s not found", metric.Name)
	case *service.CounterMetric:
		store.CMx.Lock()
		value, ok := store.Counters[metric.Name]
		store.CMx.Unlock()
		if ok {
			//metric.Value = value
			return value, nil
		}
		return 0, fmt.Errorf("storage: metric %s not found", metric.Name)
	default:
		return 0, fmt.Errorf("func (store *MemStorage) GetValue: this should never happen")
	}
}

func (store *MemStorage) GetMetric(m *service.MetricHolder) (*service.MetricHolder, error) {
	switch metric := (*m).(type) {
	case *service.GaugeMetric:
		store.GMx.Lock()
		value, ok := store.Gauges[metric.Name]
		store.GMx.Unlock()
		if ok {
			metric.Value = value
			val := service.MetricHolder(metric)
			return &val, nil
		}
		return nil, fmt.Errorf("storage: metric %s not found", metric.Name)
	case *service.CounterMetric:
		store.CMx.Lock()
		value, ok := store.Counters[metric.Name]
		store.CMx.Unlock()
		if ok {
			metric.Value = value
			val := service.MetricHolder(metric)
			return &val, nil
		}
		return nil, fmt.Errorf("storage: metric %s not found", metric.Name)
	default:
		return nil, fmt.Errorf("func (store *MemStorage) GetMetric: this should never happen")
	}
}

type StoredList struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func (store *MemStorage) MakeList() *StoredList {
	l := &StoredList{
		Gauges:   map[string]float64{},
		Counters: map[string]int64{},
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		store.GMx.Lock()
		for k, v := range store.Gauges {
			l.Gauges[k] = v
		}
		store.GMx.Unlock()
		wg.Done()
	}()
	go func() {
		store.CMx.Lock()
		for k, v := range store.Counters {
			l.Counters[k] = v
		}
		store.CMx.Unlock()
		wg.Done()
	}()
	wg.Wait()
	return l
}
