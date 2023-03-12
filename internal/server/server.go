package server

import (
	"net/http"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
)

//const ServerAddress = "localhost:8080"

const ServerAddress = "127.0.0.1:8080"

func Start(s storage.Repository) {
	r := chi.NewRouter()

	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/", handlers.GetAllMetricsHandler(s))
	http.ListenAndServe(":8080", r)
}
