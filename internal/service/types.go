package service

import (
	"fmt"
)

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
