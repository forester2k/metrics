package main

import (
	"fmt"
	"github.com/forester2k/metrics/internal/handlers"
	"github.com/forester2k/metrics/internal/service"
	"github.com/forester2k/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
	"time"
)

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
	fmt.Println("Running server on", flagRunAddr)
	mux := chi.NewRouter()
	mux.Get("/", handlers.ListStoredHandler)
	mux.Get("/value/{mType}/{mName}", handlers.ReadStoredHandler)
	mux.Post("/update/{mType}/{mName}/{mValue}", handlers.Webhook)
	return http.ListenAndServe(flagRunAddr, mux)
}
