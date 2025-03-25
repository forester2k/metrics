package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/forester2k/metrics/internal/logger"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

// HandleFile Validates the given file path. Depending on the presence of the file
// and the parameter [isRestore], it starts the function
// [ *MemStorage.readStoredData()] of reading data from this file
// or creates a file at the specified path for further saving data in it.
func HandleFile(path string, isRestore bool) (string, error) {
	path = filepath.Clean(path)
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("HandleFile: Ошибка при получении абсолютного пути %s: %w", path, err)
	}
	var isFileExists bool
	if _, err = os.Stat(path); err == nil {
		isFileExists = true
	}
	if !isFileExists && isRestore {
		return "", fmt.Errorf("HandleFile: wrong file name for load: %w", err)
	}
	if isFileExists && isRestore {
		err := Store.readStoredData(path)
		if err != nil {
			return "", fmt.Errorf("HandleFile: can't read stored data %w", err)
		}
		return path, nil
	}
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return "", fmt.Errorf("HandleFile: can't make dir %s: %w", dir, err)
	}
	var file *os.File
	file, err = os.Create(path)
	if err != nil {
		return "", fmt.Errorf("HandleFile: can't create file: %s, %w", path, err)
	}
	_ = file.Close()
	return path, nil
}

// Reads and deserializes data from a file with a given path.
func (store *MemStorage) readStoredData(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("readStoredData: can't read file, %w", err)
	}
	store.GMx.Lock()
	store.CMx.Lock()
	err = json.Unmarshal(data, store)
	store.GMx.Unlock()
	store.CMx.Unlock()
	if err != nil {
		return fmt.Errorf("readStoredData: can't unmarshal data, %w", err)
	}
	return nil
}

// Serializes and writes the current metrics to a file at the given path.
func (store *MemStorage) WriteStoreFile(path string) error {
	store.GMx.Lock()
	store.CMx.Lock()
	data, err := json.MarshalIndent(store, "", "   ")
	store.GMx.Unlock()
	store.CMx.Unlock()
	if err != nil {
		return fmt.Errorf("WriteStoreFile: can't marshal data, %w", err)
	}
	return os.WriteFile(path, data, 0666)
}

// FileStorageHandler depending on the parameter [flagStoreInterval] runs different scenarios
// for saving data to a file.
//   - [ storeBySinqSave(storePath) ] for flagStoreInterval=0
//   - [ storeByTicker(storePath, flagStoreInterval) ] for any other flagStoreInterval
func FileStorageHandler(ctx context.Context, storePath string, flagStoreInterval uint64) {
	// тут будем отслеживать время и запускать сохранение файла через интервал
	if flagStoreInterval == 0 {
		storeBySinqSave(ctx, storePath)
	} else {
		storeByTicker(ctx, storePath, flagStoreInterval)
	}
}

// Tracks time intervals [flagStoreInterval] and runs the function
// [*MemStorage.WriteStoreFile()] to save current metrics to a file
func storeByTicker(ctx context.Context, storePath string, flagStoreInterval uint64) {
	storeTicker := time.NewTicker(time.Duration(flagStoreInterval) * time.Second)
	defer storeTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-StoreSynqSave:
			// Do noting
		case <-storeTicker.C:
			err := Store.WriteStoreFile(storePath)
			if err != nil {
				logger.Log.Error("storeByTicker: can't write file", zap.Error(err))
			}
		}
	}
}

// Monitors the channel [ StoreSynqSave ] and runs the function
// [*MemStorage.WriteStoreFile()] to save current metrics to a file
func storeBySinqSave(ctx context.Context, storePath string) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-StoreSynqSave:
			err := Store.WriteStoreFile(storePath)
			if err != nil {
				logger.Log.Error("storeBySinqSave: can't write file", zap.Error(err))
			}
		}
	}
}
