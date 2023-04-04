package server

import (
	"log"
	"net/http"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
)

var serverAddress = "127.0.0.1:8080"

func SetAddress(addr string) {
	serverAddress = addr
}

func Start(s storage.Repository) {

	r := chi.NewRouter()

	handler := handlers.CompressHandle(r)
	handler = handlers.DecompressHanlder(handler)

	server := &http.Server{
		Addr:    serverAddress,
		Handler: handler,
	}

	r.Post("/update/", handlers.AddMetricJSONHandler(s))
	r.Post("/value/", handlers.GetMetricJSONHandler(s))
	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/", handlers.GetAllMetricsHandler(s))
	log.Fatal(server.ListenAndServe())
}
