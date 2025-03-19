package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

// flagRunAddr define server host:port.
// flagStoreInterval define interval in seconds to store data to file.
// flagFileStoragePath define file path to store data.
// flagRestore if true server will read data from file in start
var (
	flagRunAddr         string
	flagStoreInterval   uint64
	flagFileStoragePath string
	flagRestore         bool
)

func parseFlags() error {
	flag.StringVar(&flagRunAddr, "a", defaultRunAddr, "address and port to run server")
	flag.Uint64Var(&flagStoreInterval, "i", defaultStoreInterval, "interval in seconds to store data to file")
	flag.StringVar(&flagFileStoragePath, "f", defaultFileStoragePath, "file path to store data")
	flag.BoolVar(&flagRestore, "r", defaultRestore, "if true server will read data from file in start")
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
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		flagStoreInterval, err = strconv.ParseUint(envStoreInterval, 10, 64)
		if err != nil {
			return fmt.Errorf("parseFlags: can't parse flagStoreInterval: %w", err)
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		flagFileStoragePath = envFileStoragePath
	}
	if flagFileStoragePath, err = filepath.Abs(flagFileStoragePath); err != nil {
		return fmt.Errorf("parseFlags: can't parse flagFileStoragePath in absolute path %w", err)
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		flagRestore, err = strconv.ParseBool(envRestore)
		if err != nil {
			return fmt.Errorf("parseFlags: can't parse flagRestore: %w", err)
		}
	}
	var isFileExists bool
	if _, err = os.Stat(flagFileStoragePath); err == nil {
		isFileExists = true
	}
	// файл не существует а загружать его нужно !!!ОШИБКА
	if !isFileExists && flagRestore {
		return fmt.Errorf("parseFlags: wrong file name for load: %w", err)
	}
	// файл не существует - !!! создать файл
	if !isFileExists {
		flagFileStoragePath, err = createPathAndFile(flagFileStoragePath)
		if err != nil {
			return fmt.Errorf("parseFlags: can't create file: %w", err)
		}
	}
	return nil
}

// Валидирует строку пути к файлу, создает недостающие директории и файл.
// Возвращает строку, содержащую абсолютный путь к созданному файлу.
func createPathAndFile(path string) (string, error) {
	path = filepath.Clean(path)
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("createPathAndFile: Ошибка при получении абсолютного пути %s: %w\n", path, err)
	}
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return "", fmt.Errorf("createPathAndFile: Ошибка при создании директории %s: %w\n", dir, err)
	}
	var file *os.File
	file, err = os.Create(path)
	if err != nil {
		return "", fmt.Errorf("createPathAndFile: can't create file: %s, %w", path, err)
	}
	err = file.Close()
	if err != nil {
		return "", fmt.Errorf("createPathAndFile: can't close file: %s, %w", path, err)
	}
	return path, nil
}
