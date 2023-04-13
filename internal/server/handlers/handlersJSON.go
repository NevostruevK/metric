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

func GetMetricJSONHandler(s storage.Repository, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sM, ok := getMetricFromRequest(w, r, "", 1)
		if !ok {
			return
		}
		m := sM[0]
		rt, err := s.GetMetric(m.MType, m.ID)
		if err != nil {
			http.Error(w, "Type "+m.MType+", id "+m.ID+" not found", http.StatusNotFound)
			fmt.Println("GetMetricJSONHandler: Type " + m.MType + ", id " + m.ID + " not found")
			return
		}

		load, ok := rt.(*metrics.Metrics)
		if !ok {
			http.Error(w, "Type "+m.MType+", id "+m.ID+" is not a metric type", http.StatusNotFound)
			fmt.Println("GetMetricJSONHandler: Type " + m.MType + ", id " + m.ID + " is not a metric type")
			return
		}

		sendResponse([]metrics.Metrics{*load}, false, hashKey, w, r)
//		sendResponse(load, hashKey, w, r)
	}
}

func AddBatchMetricJSONHandler(s storage.Repository, hashKey string, cap int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
//fmt.Println("AddBatchMetricJSONHandler start")
		m, ok := getMetricFromRequest(w, r, hashKey, cap)
		if !ok {
			return
		}
fmt.Println("AddBatchMetricJSONHandler input")
fmt.Println("Size of m ",len(m))
fmt.Println(m)
		if err := s.AddGroupOfMetrics(m); err!=nil{
			fmt.Println("AddGroupOfMetrics returned the error : ",err)
			http.Error(w, fmt.Sprintf("AddGroupOfMetrics returned the error : %v",err), http.StatusInternalServerError)
			return
		}
//fmt.Println("AddBatchMetricJSONHandler stage2")

		sendResponse(m, true, hashKey, w, r)
//		fmt.Println(m)
//		fmt.Fprintln(w, m)
//		w.Write([]byte(""))
//		fmt.Fprintln(w, "Batch of Metrucs was saved normaly")
//fmt.Println("AddBatchMetricJSONHandler stage3")
	}
}

func AddMetricJSONHandler(s storage.Repository, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sM, ok := getMetricFromRequest(w, r, hashKey, 1)
		if !ok {
			return
		}
		/*		m := sM[0]
				if hashKey != "" {
					ok, err := m.CheckHash(hashKey)
					if err != nil {
						http.Error(w, "error : can't check hash for metric "+m.String(), http.StatusInternalServerError)
						fmt.Println("GetMetricJSONHandler: : can't check hash for metric " + m.String())
						return
					}
					if !ok {
						http.Error(w, "error : wrong hash for metric "+m.String(), http.StatusBadRequest)
						fmt.Println("GetMetricJSONHandler: : wrong hash for metric " + m.String())
						return
					}
				}
		*/
		if err := s.AddMetric(&sM[0]); err != nil{
			fmt.Println("AddMetric returned the error : ",err)
			http.Error(w, fmt.Sprintf("AddMetric returned the error : %v",err), http.StatusInternalServerError)
			return
		}
		sendResponse(sM, false, hashKey, w, r)
	}
}

func getMetricFromRequest(w http.ResponseWriter, r *http.Request, hashKey string, initialCapacity int) ([]metrics.Metrics, bool) {
	w.Header().Set("Content-Type", "application/json")

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "error Content-Type", http.StatusBadRequest)
		return nil, false
	}

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "read body request with an error: ", http.StatusBadRequest)
		return nil, false
	}
	sM := make([]metrics.Metrics, 0, initialCapacity)
	//	m := metrics.Metrics{}
	//err = json.Unmarshal(b, &m)
	//sM = append(sM, m)
	if initialCapacity == 1 {
		m := metrics.Metrics{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			http.Error(w, fmt.Sprintf("decode JSON with an error: %v", err), http.StatusBadRequest)
			return nil, false
		}
		sM = append(sM, m)
	} else {
		err = json.Unmarshal(b, &sM)
		if err != nil {
			http.Error(w, fmt.Sprintf("decode JSON with an error: %v", err), http.StatusBadRequest)
			return nil, false
		}
	}
	/*	if initialCapacity==1 && len(sM)!=1{
			http.Error(w, fmt.Sprintf("expected one Json metric but recieved %d metrics",len(sM)),http.StatusBadRequest)
	//		http.Error(w, "expected one Json metric but recieved "+encode.(len(sM)), http.StatusBadRequest)
			return nil, false
		}
	*/
	for _, m := range sM {
		if m.MType == "" || m.ID == "" {
			http.Error(w, "param is missed", http.StatusBadRequest)
			return nil, false
		}

		if isValidType := metrics.IsMetricType(m.MType); !isValidType {
			http.Error(w, "Type "+m.MType+" is not implemented", http.StatusNotImplemented)
			return nil, false
		}
		if hashKey != "" {
			ok, err := m.CheckHash(hashKey)
			if err != nil {
				http.Error(w, "error : can't check hash for metric "+m.String(), http.StatusInternalServerError)
				fmt.Println("GetMetricJSONHandler: : can't check hash for metric " + m.String())
				return nil, false
			}
			if !ok {
				http.Error(w, "error : wrong hash for metric "+m.String(), http.StatusBadRequest)
				fmt.Println("GetMetricJSONHandler: : wrong hash for metric " + m.String())
				return nil, false
			}
		}
	}

	return sM, true
}

func sendResponse(sM []metrics.Metrics, sendSlice bool, hashKey string, w http.ResponseWriter, r *http.Request) {

	for i, m:=range sM{
		if hashKey != "" {
			if err := sM[i].SetHash(hashKey); err != nil {
				http.Error(w, "Can't set hash for metric "+m.String(), http.StatusInternalServerError)
				return
			}
		}	
	}
	var data []byte
	var err error
	if sendSlice{
		fmt.Println("data, err = json.Marshal(&sM)")
		data, err = json.Marshal(&sM)
	}else{
		data, err = json.Marshal(sM[0])
	}	
	if err != nil {
		http.Error(w, "Can't convert to JSON", http.StatusInternalServerError)
		return
	}
	if sendSlice{
		fmt.Printf("output %d\n",len(sM))
		for _, m := range(sM){
			fmt.Println(m.String(), "  hash :",m.Hash)
		}
//		fmt.Println(sM)
	}
	w.Write(data)
}
/*
func sendResponse(m *metrics.Metrics, hashKey string, w http.ResponseWriter, r *http.Request) {

	if hashKey != "" {
		if err := m.SetHash(hashKey); err != nil {
			http.Error(w, "Can't set hash", http.StatusInternalServerError)
			return
		}
	}

	data, err := json.Marshal(*m)
	if err != nil {
		http.Error(w, "Can't convert to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
*/