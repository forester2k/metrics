package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
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
func parseFlags() error {
	flag.StringVar(&flagRunAddr, "a", defaultHost, "address and port for connecting to server")
	flag.Uint64Var(&flagReportInterval, "r", defaultReportInterval, "interval in seconds to send metrics to server")
	flag.Uint64Var(&flagPollInterval, "p", defaultPollInterval, "interval in seconds to collect metrics")
	flag.Parse()
	unknownArgs := flag.Args()
	if len(unknownArgs) > 0 {
		return fmt.Errorf("parseFlags: unknow flag %s", unknownArgs)
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	var err error
	_, err = url.ParseRequestURI(flagRunAddr)
	if err != nil {
		return fmt.Errorf("parseFlags: ParseRequestURI %w", err)
	}
	_, port, err := net.SplitHostPort(flagRunAddr)
	if err != nil {
		return fmt.Errorf("parseFlags: SplitHostPort %w", err)
	}
	_, err = strconv.ParseUint(port, 10, 16)
	if err != nil {
		return fmt.Errorf("parseFlags: incorrect port, %w", err)
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		flagReportInterval, err = strconv.ParseUint(envReportInterval, 10, 64)
		if err != nil {
			return fmt.Errorf("parseFlags: can't parse flagReportInterval: %w", err)
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		flagPollInterval, err = strconv.ParseUint(envPollInterval, 10, 64)
		if err != nil {
			return fmt.Errorf("parseFlags: can't parse flagPollInterval: %w", err)
		}
	}
	return nil
}
