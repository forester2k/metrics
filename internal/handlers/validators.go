package handlers

import (
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"strconv"
)

func URLValidate(u []string) (service.Saver, error) {
	var res service.Saver
	if len(u) <= 2 {
		return nil, fmt.Errorf("NotFound")
	}
	switch {
	case u[0] != "update" && u[0] != "value":
		return nil, fmt.Errorf("BadRequest")
	case u[1] != "gauge" && u[1] != "counter":
		return nil, fmt.Errorf("BadRequest")
	}
	switch u[1] {
	case "gauge":
		res = service.GaugeMetric{Name: u[2]}
	case "counter":
		res = service.CounterMetric{Name: u[2]}
	}
	return res, nil
}

func ValueValidate(u []string, m service.Saver) (service.Saver, error) {
	switch res := m.(type) {
	case service.GaugeMetric:
		v, err := strconv.ParseFloat(u[3], 64)
		if err != nil {
			return nil, fmt.Errorf("BadRequest")
		}
		res.Value = v
		return res, nil
	case service.CounterMetric:
		v, err := strconv.ParseInt(u[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("BadRequest")
		}
		res.Value = v
		return res, nil
	default:
		return nil, fmt.Errorf("this should never happen")
	}
}
