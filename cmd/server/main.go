package main

import (
	"fmt"
	"github.com/forester2k/metrics/internal/handlers"
	"github.com/forester2k/metrics/internal/service"
	"github.com/forester2k/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

//var Store *storage.MemStorage

// var R *chi.Mux

func main() {
	parseFlags()

	_ = service.GaugeMetric{}

	storage.Store = &storage.MemStorage{}
	storage.Store.Init()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Running server on", flagRunAddr)
	//mux := service.Mux
	mux := chi.NewRouter()

	mux.Get("/", handlers.ListStoredHandler)
	mux.Post("/{mUpdate}", handlers.Webhook)
	mux.Post("/{mUpdate}/{mType}", handlers.Webhook)
	mux.Get("/{mUpdate}/{mType}/{mName}", handlers.ReadStoredHandler)
	mux.Post("/{mUpdate}/{mType}/{mName}/{mValue}", handlers.Webhook)

	return http.ListenAndServe(flagRunAddr, mux)
}
