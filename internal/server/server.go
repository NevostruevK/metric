package server

import (
	"log"
	"net/http"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
)

//const ServerAddress = "localhost:8080"
const ServerAddress = "127.0.0.1:8080"



func Start(s storage.Repository){
		server := &http.Server{
		Addr: ServerAddress,
	}
	
	http.HandleFunc("/update/", handlers.AddMetricHandler(s))
	http.HandleFunc("/", http.NotFound)

	log.Fatal(server.ListenAndServe())
}