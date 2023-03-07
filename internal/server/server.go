package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NevostruevK/metric/internal/util/metrics"
)

const ServerAddress = "http://localhost:8080/"
//const ServerAddress = "http://127.0.0.1:8080/"


type SaveStorage struct{
	Data map[string]metrics.Metric
}
func (h SaveStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        m := metrics.NewCounterMetric("testCounterMetric",123)
        h.Data[m.Name()]=*m;
		fmt.Println("URL : ",r.URL)
		fmt.Println("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Server response")
} 


func Start(saveStorageHandler http.Handler){
	server := &http.Server{
		Addr: "localhost:8080",
	}
	

	http.Handle("/", saveStorageHandler)

	log.Fatal(server.ListenAndServe())
//	log.Fatal(http.ListenAndServe("localhost:8000", nil))
//	log.Fatal(http.ListenAndServe(":8080", nil))
}