package server

import (
	"log"
	"net/http"
	"os"

	"github.com/NevostruevK/metric/internal/db"
	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
)

const initialBatchMetricCapacity = 200

func Start(s storage.Repository, db *db.DB, address, hashKey string) {

	logger := log.New(os.Stdout, `server: `, log.LstdFlags)
	r := chi.NewRouter()

	handler := handlers.CompressHandle(r)
	handler = handlers.DecompressHanlder(handler)
	handler = handlers.LoggerHanlder(handler, logger)

	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	r.Post("/updates/", handlers.AddBatchMetricJSONHandler(s, hashKey, initialBatchMetricCapacity))
	r.Post("/update/", handlers.AddMetricJSONHandler(s, hashKey))
	r.Post("/value/", handlers.GetMetricJSONHandler(s, hashKey))
	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/ping", handlers.GetPingHandler(db))
	r.Get("/", handlers.GetAllMetricsHandler(s))
	logger.Fatal(server.ListenAndServe())
}
