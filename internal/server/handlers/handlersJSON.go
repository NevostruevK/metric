package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

func getMetric(w http.ResponseWriter, r *http.Request) (*metrics.Metrics, bool){
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json"{
		http.Error(w, "error Content-Type", http.StatusBadRequest)
		return nil, false
	}

	b, err := io.ReadAll(r.Body) 
	defer r.Body.Close()
	if err != nil{
		http.Error(w, "read body request with an error: ", http.StatusBadRequest)
		return nil, false
	}	

	m := metrics.Metrics{}
	err = json.Unmarshal(b, &m) 
	if err != nil{
		http.Error(w, "decode JSON with an error: ", http.StatusBadRequest)
		return nil, false
	}	

	if m.MType == "" || m.ID == ""{
		http.Error(w, "param is missed", http.StatusBadRequest)
		return nil, false
	}

	if isValidType := metrics.IsMetricType(m.MType); !isValidType {
		http.Error(w, "Type "+m.MType+" is not implemented", http.StatusNotImplemented)
		return nil, false
	}
	return &m, true
}

func GetMetricJSONHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m, ok:= getMetric(w, r) 
		if !ok{
			return
		}
//		fmt.Println("Get metric value",m)
		rt, err := s.GetMetric(m.MType, m.ID)
		if err != nil {
			http.Error(w, "Type "+m.MType+", id "+m.ID+" not found", http.StatusNotFound)
			fmt.Println("GetMetricJSONHandler: Type "+m.MType+", id "+m.ID+" not found")
			return
		}
		load, ok := rt.(metrics.Metrics)
		if !ok{
			http.Error(w, "Type "+m.MType+", id "+m.ID+" is not a metric type", http.StatusNotFound)
			fmt.Println("GetMetricJSONHandler: Type "+m.MType+", id "+m.ID+" is not a metric type")
			return
		}

		data, err := json.Marshal(load)
		if err != nil {
			http.Error(w, "Can't convert to JSON", http.StatusInternalServerError)
			return
		}
	
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
//		fmt.Println(data)
		w.Write(data)
	}
}

func AddMetricJSONHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m, ok:= getMetric(w, r) 
		if !ok{
			return
		}
//		fmt.Println("Add metric value",*m)

		s.AddMetric(*m)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
//		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, m)
	}
}

func ListenPOSTDefaultHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		fmt.Println("ListenPOSTDefaultHandle : ",url)
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "ListenPOSTDefaultHandle : "+url, http.StatusNotFound)
	}
}
