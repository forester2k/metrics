package main

import (
	"flag"
	"fmt"
	"os"
)

var flagRunAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	unknownArgs := flag.Args()
	if len(unknownArgs) > 0 {
		errMsg := fmt.Sprintf("unknow flag %s", unknownArgs)
		panic(errMsg)
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
}
