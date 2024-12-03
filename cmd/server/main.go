package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type saver interface {
	save(*MemStorage) error
}

type gaugeMetric struct {
	name  string
	value float64
}

func (g gaugeMetric) save(*MemStorage) error {
	// ToDo
	return nil
}

type counterMetric struct {
	name  string
	value int64
}

func (g counterMetric) save(*MemStorage) error {
	// ToDo
	return nil
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func urlValidate(u []string) (saver, error) {
	var res saver
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
		res = gaugeMetric{u[2], v}
	case "counter":
		v, err := strconv.ParseInt(u[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("BadRequest")
		}
		res = counterMetric{u[2], v}
	}

	return res, nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlSlice := strings.Split(strings.Trim(r.URL.Path, "/ "), "/")
	newMetric, err := urlValidate(urlSlice)
	if err != nil {
		switch err.Error() {
		case "BadRequest":
			w.WriteHeader(http.StatusBadRequest)
		case "NotFound":
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}

	// Отладочный блок...
	w.Header().Set("Content-Type", "text/plain")
	errMsg := "Ошибок нет"
	//if err != nil {
	//	errMsg = err.Error()
	//}
	rr := fmt.Sprintf("%#v, %v", newMetric, errMsg)

	_, _ = w.Write([]byte(rr))
}

