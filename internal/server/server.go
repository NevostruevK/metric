package server

import (
	"log"
	"net/http"

	"github.com/NevostruevK/metric/internal/db"
	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Start(s storage.Repository, db *db.DB, address, hashKey string) {

	r := chi.NewRouter()

	handler := handlers.CompressHandle(r)
	handler = handlers.DecompressHanlder(handler)

	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	r.Post("/update/", handlers.AddMetricJSONHandler(s, hashKey))
	r.Post("/value/", handlers.GetMetricJSONHandler(s, hashKey))
	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/ping", handlers.GetPingHandler(db))
	r.Get("/", handlers.GetAllMetricsHandler(s))
	log.Fatal(server.ListenAndServe())
}
