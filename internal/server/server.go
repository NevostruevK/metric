package server

import (
	"net/http"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
)

//const ServerAddress = "localhost:8080"
const ServerAddress = "127.0.0.1:8080"



func Start(s storage.Repository){
	r := chi.NewRouter()

//	r.Use(middleware.RequestID)
//    r.Use(middleware.RealIP)
//    r.Use(middleware.Logger)
//    r.Use(middleware.Recoverer)
	
	r.Post("/update/{typeM}/{name}/{value}",handlers.AddMetricHandler(s))	
	r.Get("/value/{typeM}/{name}",handlers.GetMetricHandler(s))	
	r.Get("/",handlers.GetAllMetricsHandler(s))	
//	r.Get("/update/{typeM}/{name}/{value}",handlers.AddMetricHandler(s))	
	http.ListenAndServe(":8080", r)

//	http.HandleFunc("/update/", handlers.AddMetricHandler(s))
//	http.HandleFunc("/", http.NotFound)

//	log.Fatal(server.ListenAndServe())
}

/*r := chi.NewRouter()
r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("chi"))
})
http.ListenAndServe(":8080", r)
*/