package handlers

import (
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"strconv"
)

func URLValidate(mType string, mName string) (service.MetricHolder, error) {
	var res service.MetricHolder
	switch mType {
	case "gauge":
		res = &service.GaugeMetric{Name: mName}
	case "counter":
		res = &service.CounterMetric{Name: mName}
	default:
		return nil, fmt.Errorf("BadRequest")
	}
	return res, nil
}

func ValueValidate(mValue string, m service.MetricHolder) (service.MetricHolder, error) {
	switch res := m.(type) {
	case *service.GaugeMetric:
		v, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return nil, fmt.Errorf("BadRequest")
		}
		res.Value = v
		return res, nil
	case *service.CounterMetric:
		v, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("BadRequest")
		}
		res.Value = v
		return res, nil
	default:
		return nil, fmt.Errorf("func ValueValidate: this should never happen")
	}
}

func JSONTypeAndValueValidate(req *service.Metrics) (service.MetricHolder, error) {
	var res service.MetricHolder
	switch req.MType {
	case "gauge":
		if req.Value == nil {
			return nil, fmt.Errorf("BadRequest")
		}
		res = &service.GaugeMetric{Name: req.ID, Value: *req.Value}
	case "counter":
		if req.Delta == nil {
			return nil, fmt.Errorf("BadRequest")
		}
		res = &service.CounterMetric{Name: req.ID, Value: *req.Delta}
	default:
		return nil, fmt.Errorf("BadRequest")
	}
	return res, nil
}

func JSONTypeValidate(req *service.Metrics) (service.MetricHolder, error) {
	var res service.MetricHolder
	switch req.MType {
	case "gauge":
		res = &service.GaugeMetric{Name: req.ID}
	case "counter":
		res = &service.CounterMetric{Name: req.ID}
	default:
		return nil, fmt.Errorf("BadRequest")
	}
	return res, nil
}
