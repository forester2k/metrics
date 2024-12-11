package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// flagRunAddr define servers host:port
// flagReportInterval define time interval in seconds between servers reports
// flagPollInterval define time interval in seconds between metrics polling
var (
	flagRunAddr        string
	flagReportInterval uint64
	flagPollInterval   uint64
)

// Define and parces flags and environment variables
func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", defaultHost, "address and port for connecting to server")
	flag.Uint64Var(&flagReportInterval, "r", defaultReportInterval, "interval in seconds to send metrics to server")
	flag.Uint64Var(&flagPollInterval, "p", defaultPollInterval, "interval in seconds to collect metrics")
	flag.Parse()
	unknownArgs := flag.Args()
	if len(unknownArgs) > 0 {
		errMsg := fmt.Sprintf("unknow flag %s", unknownArgs)
		panic(errMsg)
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	var err error
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		flagReportInterval, err = strconv.ParseUint(envReportInterval, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		flagPollInterval, err = strconv.ParseUint(envPollInterval, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}
}
