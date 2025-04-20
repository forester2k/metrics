package main

import (
	"github.com/forester2k/metrics/internal/logger"
	"github.com/forester2k/metrics/internal/service"
	"go.uber.org/zap"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var m runtime.MemStats
var random *rand.Rand
var specialMetrics = map[string]service.MetricHolder{
	"RandomValue": &service.GaugeMetric{Name: "RandomValue", Value: float64(0)},
	"PollCount":   &service.CounterMetric{Name: "PollCount", Value: int64(0)},
}

func pollingInit() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func poll(mutex *sync.Mutex) {
	mutex.Lock()
	runtime.ReadMemStats(&m)
	err := specialMetrics["RandomValue"].SetValue(random.Float64())
	if err != nil {
		logger.Log.Error("poll: error to set RandomValue", zap.Error(err))
	}
	val := specialMetrics["PollCount"].GetValue().(int64)
	val++
	err = specialMetrics["PollCount"].SetValue(val)
	if err != nil {
		logger.Log.Error("poll: error to set PollCount", zap.Error(err))
	}
	mutex.Unlock()
}
