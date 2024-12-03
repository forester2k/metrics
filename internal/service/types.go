package service

import "fmt"

type Saver interface {
	Save(*MemStorage) error
	Path() string
}

type PathMaker interface {
	PathMake() string
}

type GaugeMetric struct {
	Name  string
	Value float64
}

func (g GaugeMetric) Save(*MemStorage) error {
	// ToDo
	return nil
}

func (m GaugeMetric) Path() string {
	return "/gauge/" + m.Name + "/" + fmt.Sprint(m.Value)
}

type CounterMetric struct {
	Name  string
	Value int64
}

func (m CounterMetric) Path() string {
	return "/counter/" + m.Name + "/" + fmt.Sprint(m.Value)
}

func (g CounterMetric) Save(*MemStorage) error {
	// ToDo
	return nil
}

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}
