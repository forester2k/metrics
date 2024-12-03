package main

import (
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"net/http"
	"time"
)

const reportInterval = 10 * time.Second
const pollInterval = 2 * time.Second
const host = "http://localhost:8080/"

func poll() {
	//fmt.Println("polling", time.Now())
	// TODO: logic
}

func report(metric service.Saver) {
	//fmt.Println("reporting", time.Now())

	endpoint := host + "update" + metric.Path()
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "text/plain")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	resp := response
	//fmt.Println(resp)
	_ = resp
	// TODO: logic
}

func main() {
	mockGaugeMetric := service.GaugeMetric{Name: "someGaugeMetric", Value: float64(3.1)}
	mockCounterMetric := service.CounterMetric{Name: "someCounterMetric", Value: int64(3)}

	reportTicker := time.Tick(reportInterval)
	pollTicker := time.Tick(pollInterval)
	for {
		select {
		case <-reportTicker:
			go report(mockGaugeMetric)
			go report(mockCounterMetric)
		case <-pollTicker:
			go poll()
		}

	}
}
