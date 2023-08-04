// Package server инициализирует сервер для работы с метриками.
package server

import (
	"log"
	"net"
	"net/http"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/crypt"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/go-chi/chi/v5"
)

const initialBatchMetricCapacity = 200

// NewServer создание сервера на основе роутера github.com/go-chi/chi/v5.
func NewServer(s storage.Repository, cfg *commands.Config) (*http.Server, error) {
	handlers.Logger = logger.NewLogger(`server: `, log.LstdFlags)

	dcr, err := crypt.NewDecrypt(cfg.CryptoKey)
	if err != nil {
		handlers.Logger.Printf("failed to create decrypt entity %v", err)
		return nil, err
	}

	r := chi.NewRouter()

	handler := handlers.CompressHandle(r)
	handler = handlers.DecompressHanlder(handler)
	handler = handlers.DecryptHanlder(handler, dcr)
	if cfg.TrustedSubnet != "" {
		_, ipNet, err := net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			handlers.Logger.Printf("failed to parse cfg.TrustedSubnet %v", err)
		} else {
			handler = handlers.IPCheckHandler(handler, ipNet)
		}
	}
	handler = handlers.LoggerHanlder(handler, handlers.Logger)

	r.Post("/updates/", handlers.AddBatchMetricJSONHandler(s, cfg.HashKey, initialBatchMetricCapacity))
	r.Post("/update/", handlers.AddMetricJSONHandler(s, cfg.HashKey))
	r.Post("/value/", handlers.GetMetricJSONHandler(s, cfg.HashKey))
	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/ping", handlers.GetPingHandler(s))
	r.Get("/", handlers.GetAllMetricsHandler(s))

	return &http.Server{
		Addr:    cfg.Address,
		Handler: handler,
	}, nil
}
