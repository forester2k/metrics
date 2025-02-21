package main

import (
	"fmt"
	"github.com/forester2k/metrics/internal/handlers"
	"github.com/forester2k/metrics/internal/logger"
	"github.com/forester2k/metrics/internal/middleware"
	"github.com/forester2k/metrics/internal/service"
	"github.com/forester2k/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const defaultLogLevel = "Info"

//const defaultLogLevel = "Debug"

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
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Printf("Running server on \"%s\"\n", flagRunAddr)
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
	return http.ListenAndServe(flagRunAddr, mux)
}
