package server

import (
	"log"
	"net/http"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
)

const ServerAddress = "127.0.0.1:8080"

func Start(s storage.Repository) {
	r := chi.NewRouter()

	r.Post("/update", handlers.AddMetricJSONHandler(s))
	r.Post("/value", handlers.GetMetricJSONHandler(s))
	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))	
	r.Post("/", handlers.ListenPOSTDefaultHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/", handlers.GetAllMetricsHandler(s))
	log.Fatal(http.ListenAndServe(":8080", r))
}
