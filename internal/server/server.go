// Package server инициализирует сервер для работы с метриками.
package server

import (
	"log"
	"net/http"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/crypt"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/go-chi/chi/v5"
)

const initialBatchMetricCapacity = 200

// Server HTTP сервер для работы с метриками.
type Server *http.Server

// NewServer создание сервера на основе роутера github.com/go-chi/chi/v5.
func NewServer(s storage.Repository, address, hashKey, cryptKey string) Server {
	handlers.Logger = logger.NewLogger(`server: `, log.LstdFlags)

	dcr, err := crypt.NewDecrypt(cryptKey)
	if err != nil {
		handlers.Logger.Printf("failed to create decrypt entity %v", err)
	}

	r := chi.NewRouter()

	handler := handlers.CompressHandle(r)
	handler = handlers.DecompressHanlder(handler)
	handler = handlers.DecryptHanlder(handler, dcr)
	handler = handlers.LoggerHanlder(handler, handlers.Logger)

	r.Post("/updates/", handlers.AddBatchMetricJSONHandler(s, hashKey, initialBatchMetricCapacity))
	r.Post("/update/", handlers.AddMetricJSONHandler(s, hashKey))
	r.Post("/value/", handlers.GetMetricJSONHandler(s, hashKey))
	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/ping", handlers.GetPingHandler(s))
	r.Get("/", handlers.GetAllMetricsHandler(s))

	return &http.Server{
		Addr:    address,
		Handler: handler,
	}
}
