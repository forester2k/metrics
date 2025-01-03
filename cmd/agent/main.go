package main

import (
	"log"
	"sync"
	"time"
)

const defaultReportInterval = 10
const defaultPollInterval = 2
const defaultHost = "localhost:8080"

func main() {
	pollingInit()
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
