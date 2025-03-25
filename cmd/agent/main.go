package main

import (
	"compress/gzip"
	"github.com/forester2k/metrics/internal/logger"
	"log"
	"sync"
	"time"
)

const (
	defaultReportInterval uint64 = 10
	defaultPollInterval   uint64 = 2
	defaultHost           string = "localhost:8080"
	gzipCompressLevel     int    = gzip.BestSpeed
	defaultLogLevel       string = "Info"
)

func main() {
	pollingInit()
	if err := logger.Initialize(defaultLogLevel); err != nil {
		log.Fatal(err)
	}
	err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	reportTicker := time.Tick(time.Duration(flagReportInterval) * time.Second)
	pollTicker := time.Tick(time.Duration(flagPollInterval) * time.Second)
	var mutex sync.Mutex
	for {
		select {
		case <-reportTicker:
			go report(&mutex)
		case <-pollTicker:
			go poll(&mutex)
		}
	}
}
