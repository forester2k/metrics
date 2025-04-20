package service

import (
	"fmt"
)

// Metrics - структура для обмена метриками в формате JSON.
//
// ID - имя метрики.
// MType - параметр, принимающий значение gauge или counter.
// Delta - значение метрики в случае передачи counter.
// Value - значение метрики в случае передачи gauge.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func ConvToMetrics(m *MetricHolder) Metrics {
	switch mTemp := (*m).(type) {
	case *GaugeMetric:
		return Metrics{ID: mTemp.Name, MType: "gauge", Value: &mTemp.Value}
	case *CounterMetric:
		return Metrics{ID: mTemp.Name, MType: "counter", Delta: &mTemp.Value}
	default:
		return Metrics{} //this should never happen
	}
}

type Storekeeper interface {
	Save(MetricHolder) error
}

type MetricHolder interface {
	GetPath() string
	GetValue() interface{}
	SetValue(interface{}) error
}

type GaugeMetric struct {
	Name  string
	Value float64
}

func (m *GaugeMetric) GetPath() string {
	return "/gauge/" + m.Name + "/" + fmt.Sprint(m.Value)
}

func (m *GaugeMetric) GetValue() interface{} {
	return m.Value
}

func (m *GaugeMetric) SetValue(i interface{}) error {
	switch n := i.(type) {
	case float64:
		m.Value = n
		return nil
	default:
		return fmt.Errorf("cant convert %v to float64", i)
	}
}

type CounterMetric struct {
	Name  string
	Value int64
}

func (m *CounterMetric) GetPath() string {
	return "/counter/" + m.Name + "/" + fmt.Sprint(m.Value)
}

func (m *CounterMetric) GetValue() interface{} {
	return m.Value
}

func (m *CounterMetric) SetValue(i interface{}) error {
	switch n := i.(type) {
	case int64:
		m.Value = n
		return nil
	default:
		return fmt.Errorf("cant convert %v to int64", i)
	}
}
