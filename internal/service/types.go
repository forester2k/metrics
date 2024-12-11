package service

import (
	"fmt"
	"github.com/forester2k/metrics/internal/storage"
)

type Saver interface {
	Save(*storage.MemStorage) error
	Get(*storage.MemStorage) (Saver, error)
	Path() string
}

type GaugeMetric struct {
	Name  string
	Value float64
}

func (m GaugeMetric) Save(store *storage.MemStorage) error {
	err := store.AddGauge(m.Name, m.Value)
	if err != nil {
		return fmt.Errorf("GaugeMetricSave: %s", err)
	}
	return nil
}

func (m GaugeMetric) Get(store *storage.MemStorage) (Saver, error) {
	value, err := store.GetGauge(m.Name)
	if err != nil {
		return m, fmt.Errorf("GaugeMetricGet: %s", err)
	}
	m.Value = value
	return m, nil
}

func (m GaugeMetric) Path() string {
	return "/gauge/" + m.Name + "/" + fmt.Sprint(m.Value)
}

type CounterMetric struct {
	Name  string
	Value int64
}

func (m CounterMetric) Save(store *storage.MemStorage) error {
	err := store.AddCounter(m.Name, m.Value)
	if err != nil {
		return fmt.Errorf("CounterMetricSave: %s", err)
	}
	return nil
}

func (m CounterMetric) Get(store *storage.MemStorage) (Saver, error) {
	value, err := store.GetCounter(m.Name)
	if err != nil {
		return m, fmt.Errorf("CounterMetricGet: %s", err)
	}
	m.Value = value
	return m, nil
}

func (m CounterMetric) Path() string {
	return "/counter/" + m.Name + "/" + fmt.Sprint(m.Value)
}

type ListTaker interface {
	ListTake(*storage.MemStorage) (MetricList, error)
}

type MetricList storage.StoredList

func (l *MetricList) TakeList(store *storage.MemStorage) {
	ll := store.MakeList()
	l.Gauges = ll.Gauges
	l.Counters = ll.Counters
}
