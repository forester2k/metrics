package storage

import (
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"sync"
)

type MemStorage struct {
	GMx      sync.Mutex
	CMx      sync.Mutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func (store *MemStorage) Save(m *service.MetricHolder) error {
	switch metric := (*m).(type) {
	case *service.GaugeMetric:
		store.GMx.Lock()
		store.Gauges[metric.Name] = metric.Value
		store.GMx.Unlock()
		return nil
	case *service.CounterMetric:
		store.CMx.Lock()
		if _, ok := store.Counters[metric.Name]; !ok {
			store.Counters[metric.Name] = 0
		}
		store.Counters[metric.Name] = store.Counters[metric.Name] + metric.Value
		store.CMx.Unlock()
		return nil
	default:
		return fmt.Errorf("func (store *MemStorage) Save: this should never happen")
	}
}

var Store *MemStorage

func (store *MemStorage) Init() {
	store.Gauges = make(map[string]float64)
	store.Counters = make(map[string]int64)
	store.Counters["PollCount"] = 0
}

func (store *MemStorage) Get(m *service.MetricHolder) (any, error) {
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
		return 0, fmt.Errorf("func (store *MemStorage) Get: this should never happen")
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
