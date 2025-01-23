package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"net/http"
	"reflect"
	"runtime"
	"sync"
)

func reportMetric(metric service.MetricHolder) error {
	endpoint := "http://" + flagRunAddr + "/update" + metric.GetPath()
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("reportMetric: can't make request: %w", err)
	}
	request.Header.Add("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("reportMetric: ошибка в client.Do(request): %w", err)
	}
	response.Body.Close()
	return nil
}

func reportJSONMetric(metric service.MetricHolder) error {
	endpoint := "http://" + flagRunAddr + "/update"
	m := service.ConvToMetrics(&metric)
	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("reportJSONMetric: can't marshal metric: %w", err)
	}
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("reportMetric: can't make request: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("reportMetric: ошибка в client.Do(request): %w", err)
	}
	response.Body.Close()
	return nil
}

func makeGaugeMetrics(m runtime.MemStats) []*service.GaugeMetric {
	gaugeList := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
		"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups",
		"MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

	res := make([]*service.GaugeMetric, 0, len(gaugeList))
	var mutex sync.Mutex
	var wg sync.WaitGroup
	v := reflect.ValueOf(m)
	for _, metricName := range gaugeList {
		wg.Add(1)
		go func(metricName string, mutex *sync.Mutex) {
			defer wg.Done()
			value := v.FieldByName(metricName)
			var val float64
			switch value.Kind() {
			case reflect.Float64:
				val = value.Float()
			case reflect.Uint64:
				val = float64(value.Uint())
			case reflect.Uint32:
				val = float64(value.Uint())
			default:
				err := fmt.Errorf("report(): не найден обработчик для типа %v", value.Kind())
				fmt.Println(err)
				val = float64(0)
			}
			mutex.Lock()
			res = append(res, &service.GaugeMetric{Name: metricName, Value: val})
			mutex.Unlock()
		}(metricName, &mutex)
	}
	wg.Wait()
	return res
}

func report(mutex *sync.Mutex) {
	mutex.Lock()
	gaugeMetrics := makeGaugeMetrics(m)
	for _, metric := range gaugeMetrics {
		err := reportJSONMetric(metric)
		if err != nil {
			// TODO вывести ошибку в лог уровня ERROR
			err = fmt.Errorf("agent report(): %w", err)
			fmt.Println(err)
		}
	}
	for _, metric := range specialMetrics {
		err := reportJSONMetric(metric)
		if err != nil {
			// TODO вывести ошибку в лог уровня ERROR
			err = fmt.Errorf("agent report(): %w", err)
			fmt.Println(err)
		}
	}
	err := specialMetrics["PollCount"].SetValue(int64(0))
	if err != nil {
		// TODO вывести ошибку в лог уровня ERROR
		fmt.Println(err)
	}
	mutex.Unlock()
}
