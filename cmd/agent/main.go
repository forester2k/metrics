package main

import (
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"net/http"
	"time"
)

const defaultReportInterval = 10
const defaultPollInterval = 2
const defaultHost = "localhost:8080"

func poll() {
	//fmt.Println("polling", time.Now())
	// TODO: logic
}

func report(metric service.Saver) {
	//fmt.Println("reporting", time.Now())

	endpoint := "http://" + flagRunAddr + "/update" + metric.Path()
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "text/plain")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка в client.Do(request)", err)
	}
	resp := response
	//fmt.Println(resp)
	_ = resp
	resp.Body.Close()
	// TODO: logic...
}

func main() {
	parseFlags()
	mockGaugeMetric := service.GaugeMetric{Name: "someGaugeMetric", Value: float64(3.1)}
	mockCounterMetric := service.CounterMetric{Name: "someCounterMetric", Value: int64(3)}
	time.Sleep(2 * time.Second)
	reportTicker := time.Tick(time.Duration(flagReportInterval) * time.Second)
	pollTicker := time.Tick(time.Duration(flagPollInterval) * time.Second)
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
