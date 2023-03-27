package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/fgzip"
	"github.com/NevostruevK/metric/internal/util/metrics"
)
func GetMetricJSONHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m, ok:= getRequest(w, r) 
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

		m, ok:= getRequest(w, r) 
		if !ok{
			return
		}
		s.AddMetric(m)

		sendResponse(m, w, r)
	}
}

func getRequest(w http.ResponseWriter, r *http.Request) (*metrics.Metrics, bool){
	w.Header().Set("Content-Type", "application/json")

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
//		if r.Header.Get("Content-Type") != "application/json"{
		http.Error(w, "error Content-Type", http.StatusBadRequest)
		return nil, false
	}

	b, err := io.ReadAll(r.Body) 
	defer r.Body.Close()
	if err != nil{
		http.Error(w, "read body request with an error: ", http.StatusBadRequest)
		return nil, false
	}	
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
//		if r.Header.Get("Content-Encoding") != "gzip"{
//		fmt.Println("find Content-Encoding gzip")
		b, err = fgzip.Decompress(b)
		if err != nil{
			http.Error(w, "can't decompress data: ", http.StatusBadRequest)
			return nil, false		
		} 
	}
//	fmt.Println("can't find Countent Encoding gzip")
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
        if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			data, err = fgzip.Compress(data)
			if err != nil{
				fmt.Println("can't compress data", err)
			}else{
				w.Header().Add("Content-Encoding", "gzip")
			}
		}	
		w.Write(data)
}		
