package handlers

import (
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"strconv"
)

func UrlValidate(u []string) (service.Saver, error) {
	var res service.Saver
	if len(u) < 4 {
		return nil, fmt.Errorf("NotFound")
	}

	switch {
	case u[0] != "update":
		return nil, fmt.Errorf("BadRequest")
	case u[1] != "gauge" && u[1] != "counter":
		return nil, fmt.Errorf("BadRequest")
	}

	switch u[1] {
	case "gauge":
		v, err := strconv.ParseFloat(u[3], 64)
		if err != nil {
			return nil, fmt.Errorf("BadRequest")
		}
		res = service.GaugeMetric{u[2], v}
	case "counter":
		v, err := strconv.ParseInt(u[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("BadRequest")
		}
		res = service.CounterMetric{u[2], v}
	}

	return res, nil
}
