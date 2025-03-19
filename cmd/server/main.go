package main

import (
	"context"
	"fmt"
	"github.com/forester2k/metrics/internal/handlers"
	"github.com/forester2k/metrics/internal/logger"
	"github.com/forester2k/metrics/internal/middleware"
	"github.com/forester2k/metrics/internal/service"
	"github.com/forester2k/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// defaultRunAddr - default server host:port.
// defaultLogLevel - default loggin level.
// defaultStoreInterval - default interval in seconds to store data to file.
// defaultFileStoragePath - default file path to store data
// defaultRestore - if true server will read data from file in start
const (
	defaultRunAddr         string = "localhost:8080"
	defaultLogLevel        string = "Info"
	defaultStoreInterval   uint64 = 300 //сделать 300 как отлажу
	defaultFileStoragePath string = "./file_st/saved.json"
	defaultRestore         bool   = false
)

var srv http.Server
var random *rand.Rand

func Init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	_ = service.GaugeMetric{}
	storage.Store = &storage.MemStorage{}
	storage.Store.Init()
	storePath, err := storage.HandleFile(flagFileStoragePath, flagRestore)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return run()
	})
	g.Go(func() error {
		storage.FileStorageHandler(gCtx, storePath, flagStoreInterval)
		return fmt.Errorf("fileStorageHandler was stoped")
	})
	g.Go(func() error {
		<-gCtx.Done()
		fmt.Println("Shutdown server...")
		return srv.Shutdown(context.Background())
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
	fmt.Println("Saving data...")
	err = storage.Store.WriteStoreFile(storePath)
	if err != nil {
		fmt.Println(fmt.Errorf("main: can't save file, %w", err))
	}
}

func run() error {
	fmt.Printf("Running server on \"%s\"\n", flagRunAddr)
	if flagRestore {
		fmt.Println("Loading metrics from ", flagFileStoragePath)
	}
	fmt.Printf("Saving metrics in %v, with interval %v seconds\n", flagFileStoragePath, flagStoreInterval)
	if err := logger.Initialize(defaultLogLevel); err != nil {
		return fmt.Errorf("main.run: %w", err)
	}
	fmt.Printf("...with \"%s\" logging level\n", defaultLogLevel)
	mux := chi.NewRouter()
	mux.Use(middleware.RequestGzipDecompressor)
	mux.Use(middleware.ResponseGzipCompressor)
	mux.Use(logger.RequestResponseLogger)
	mux.Get("/", handlers.ListStoredHandler)
	mux.Post("/update", handlers.WriteJSONMetricHandler)
	mux.Post("/update/", handlers.WriteJSONMetricHandler)
	mux.Post("/value", handlers.ReadJSONMetricHandler)
	mux.Post("/value/", handlers.ReadJSONMetricHandler)
	mux.Get("/value/{mType}/{mName}", handlers.ReadStoredHandler)
	mux.Post("/update/{mType}/{mName}/{mValue}", handlers.Webhook)
	//return http.ListenAndServe(flagRunAddr, mux)
	srv.Addr = flagRunAddr
	srv.Handler = mux
	fmt.Println("+++ Ща как сервер запущу !!!!")
	return srv.ListenAndServe()
}
