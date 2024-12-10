package storage

import (
	"fmt"
)

type MemStorage struct {
	Gauges    map[string]float64
	Counters  map[string]int64
	ValidKeys map[string]string
}

var Store *MemStorage

func (m *MemStorage) Init() {
	m.ValidKeys = GetValidKeys()
	m.Gauges = make(map[string]float64)
	m.Counters = make(map[string]int64)
	m.Counters["PollCount"] = 0
	//m.Gauges["one"] = 1.1
	//fmt.Println("from storage", m)
}

//func isValidMetric(name string, mType string, m *MemStorage) error {
//	if _, ok := m.ValidKeys[name]; !ok {
//		return nil
//		//return fmt.Errorf("storage: %s is not valid metric name", name)
//	}
//	if m.ValidKeys[name] != mType {
//		return nil
//		//return fmt.Errorf("storage: %s is not valid metric type", m.ValidKeys[name])
//	}
//	return nil
//}

func (m *MemStorage) AddGauge(name string, value float64) error {
	//if err := isValidMetric(name, "gauge", m); err != nil {
	//	return err
	//}
	m.Gauges[name] = value
	return nil
}

func (m *MemStorage) AddCounter(name string, value int64) error {
	//if err := isValidMetric(name, "counter", m); err != nil {
	//	return err
	//}
	if _, ok := m.Counters[name]; !ok {
		m.Counters[name] = 0
	}
	m.Counters[name] = m.Counters[name] + value
	return nil
}

func (m *MemStorage) GetGauge(name string) (float64, error) {
	//if err := isValidMetric(name, "gauge", m); err != nil {
	//	return 0, err
	//}
	value, ok := m.Gauges[name]
	if ok {
		return value, nil
	}
	return 0, fmt.Errorf("storage: metric %s not found", name)
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
	//if err := isValidMetric(name, "counter", m); err != nil {
	//	return 0, err
	//}
	value, ok := m.Counters[name]
	if ok {
		return value, nil
	}
	return 0, fmt.Errorf("storage: metric %s not found", name)
}

type StoredList struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func (m *MemStorage) MakeList() *StoredList {
	l := &StoredList{
		Gauges:   map[string]float64{},
		Counters: map[string]int64{},
	}
	for k, v := range m.Gauges {
		l.Gauges[k] = v
	}
	for k, v := range m.Counters {
		l.Counters[k] = v
	}
	return l
}
