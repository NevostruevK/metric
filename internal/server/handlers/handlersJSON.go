package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
)
func GetMetricJSONHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m, ok:= getMetricFromRequest(w, r) 
		if !ok{
			return
		}

		rt, err := s.GetMetric(m.MType, m.ID)
		if err != nil {
			http.Error(w, "Type "+m.MType+", id "+m.ID+" not found", http.StatusNotFound)
			fmt.Println("GetMetricJSONHandler: Type "+m.MType+", id "+m.ID+" not found")
			return
		}

		load, ok := rt.(*metrics.Metrics)
		if !ok{
			http.Error(w, "Type "+m.MType+", id "+m.ID+" is not a metric type", http.StatusNotFound)
			fmt.Println("GetMetricJSONHandler: Type "+m.MType+", id "+m.ID+" is not a metric type")
			return
		}
		sendResponse(load, w, r)
	}
}

func AddMetricJSONHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m, ok:= getMetricFromRequest(w, r) 
		if !ok{
			return
		}
		s.AddMetric(m)

		sendResponse(m, w, r)
	}
}

func getMetricFromRequest(w http.ResponseWriter, r *http.Request) (*metrics.Metrics, bool){
	w.Header().Set("Content-Type", "application/json")

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
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

func sendResponse(m *metrics.Metrics, w http.ResponseWriter, r *http.Request){
		data, err := json.Marshal(*m)
		if err != nil {
			http.Error(w, "Can't convert to JSON", http.StatusInternalServerError)
			return
		}
		w.Write(data)
}		
