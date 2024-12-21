package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
)

var flagRunAddr string

func parseFlags() error {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
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
	return nil
}
