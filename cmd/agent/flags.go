package main

import (
	"flag"
	"fmt"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var (
	flagRunAddr        string
	flagReportInterval uint
	flagPollInterval   uint
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", defaultHost, "address and port for connecting to server")
	flag.UintVar(&flagReportInterval, "r", defaultReportInterval, "interval in seconds to send metrics to server")
	flag.UintVar(&flagPollInterval, "p", defaultPollInterval, "interval in seconds to collect metrics")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	unknownArgs := flag.Args()
	if len(unknownArgs) > 0 {
		errMsg := fmt.Sprintf("unknow flag %s", unknownArgs)
		panic(errMsg)
	}
}
