package server

import (
	"log"
	"net/http"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/go-chi/chi/v5"
)

const initialBatchMetricCapacity = 200

//	type Server struct{
//		server *http.Server
//	}
type Server *http.Server

func NewServer(s storage.Repository, address, hashKey string) Server {
	handlers.Logger = logger.NewLogger(`server: `, log.LstdFlags)
	r := chi.NewRouter()

	handler := handlers.CompressHandle(r)
	handler = handlers.DecompressHanlder(handler)
	handler = handlers.LoggerHanlder(handler, handlers.Logger)

	r.Post("/updates/", handlers.AddBatchMetricJSONHandler(s, hashKey, initialBatchMetricCapacity))
	r.Post("/update/", handlers.AddMetricJSONHandler(s, hashKey))
	r.Post("/value/", handlers.GetMetricJSONHandler(s, hashKey))
	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/ping", handlers.GetPingHandler(s))
	r.Get("/", handlers.GetAllMetricsHandler(s))

	//	return &Server{

	return &http.Server{
		Addr:    address,
		Handler: handler,
	}
	// }
}

/*
func (s Server) Start(){
	handlers.Logger.Fatal(s.ListenAndServe())
}
*/
/*
func Start(s storage.Repository, address, hashKey string) {

	handlers.Logger = logger.NewLogger(`server: `, log.LstdFlags)
	r := chi.NewRouter()

	handler := handlers.CompressHandle(r)
	handler = handlers.DecompressHanlder(handler)
	handler = handlers.LoggerHanlder(handler, handlers.Logger)

	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	r.Post("/updates/", handlers.AddBatchMetricJSONHandler(s, hashKey, initialBatchMetricCapacity))
	r.Post("/update/", handlers.AddMetricJSONHandler(s, hashKey))
	r.Post("/value/", handlers.GetMetricJSONHandler(s, hashKey))
	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))
	r.Get("/ping", handlers.GetPingHandler(s))
	r.Get("/", handlers.GetAllMetricsHandler(s))
	handlers.Logger.Fatal(server.ListenAndServe())
}
*/
