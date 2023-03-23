package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/caarlos0/env/v7"
	"github.com/go-chi/chi/v5"
)

const ServerAddress = "127.0.0.1:8080"

type environment struct{
	Address 	string 	`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func Start(s storage.Repository) {
	en := environment{}
	if err := env.Parse(&en); err != nil{
		fmt.Printf("Server read environment with the error: %+v\n", err)
	}

	r := chi.NewRouter()
	
	server := &http.Server{
		Addr: en.Address,
		Handler: r,
	}

	r.Post("/update/", handlers.AddMetricJSONHandler(s))
	r.Post("/value/", handlers.GetMetricJSONHandler(s))
	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))	
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/", handlers.GetAllMetricsHandler(s))
	log.Fatal(server.ListenAndServe())
}
