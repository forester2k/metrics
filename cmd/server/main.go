package main

import (
	"github.com/forester2k/metrics/internal/handlers"
	"github.com/forester2k/metrics/internal/service"
	"github.com/forester2k/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

//var Store *storage.MemStorage

// var R *chi.Mux

func main() {

	_ = service.GaugeMetric{}

	storage.Store = &storage.MemStorage{}
	storage.Store.Init()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := service.Mux
	mux = chi.NewRouter()

	mux.Get("/", handlers.ListStoredHandler)
	mux.Post("/{mUpdate}", handlers.Webhook)
	mux.Post("/{mUpdate}/{mType}", handlers.Webhook)
	mux.Get("/{mUpdate}/{mType}/{mName}", handlers.ReadStoredHandler)
	mux.Post("/{mUpdate}/{mType}/{mName}/{mValue}", handlers.Webhook)

	return http.ListenAndServe("localhost:8080", mux)
}
