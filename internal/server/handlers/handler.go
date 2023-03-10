package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/go-chi/chi/v5"
)

func GetAllMetricsHandler(s storage.Repository) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
	//получаем список всех пользователей
	sm := s.GetAllMetrics()
	//указываем путь к файлу с шаблоном
	main := filepath.Join(".", "getAllMetrics.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(main)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//исполняем именованный шаблон "users", передавая туда массив со списком пользователей
	err = tmpl.ExecuteTemplate(w, "Metric", sm)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}
}



func GetMetricHandler(s storage.Repository) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                typeM := chi.URLParam(r, "typeM")
                name := chi.URLParam(r, "name")
                fmt.Println("typeM :",typeM,"name :",name)
                if typeM == "" || name == "" {
                        http.Error(w, "param is missed", http.StatusBadRequest)
                        fmt.Println("type: ",typeM," name: ",name)
                        return
                }
                if validType := metrics.IsMetricType(typeM); !validType {
                        http.Error(w, "Type "+typeM+" is not implemented", http.StatusNotImplemented)
                        return
                }
                        m, err := s.GetMetric(typeM, name)
                if err != nil{
                        http.Error(w, "Type "+typeM+", Name "+name+" not found", http.StatusNotFound)
                        return
                }
                w.WriteHeader(http.StatusOK)
                w.Header().Set("Content-Type", "text/plain")
                fmt.Fprintln(w, m.Value())
        }
}

//func SaveMetricHandler(s *storage.MemStorage) http.HandlerFunc {
func AddMetricHandler(s storage.Repository) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                typeM := chi.URLParam(r, "typeM")
                name := chi.URLParam(r, "name")
                value := chi.URLParam(r, "value")
//                fmt.Println("Car ID :",typeM,"name :",name,"value :",value)
                if typeM == "" || name == "" || value == ""{
                        http.Error(w, "param is missed", http.StatusBadRequest)
                        fmt.Println("type: ",typeM," name: ",name," value: ",value)
                        return
                }
                var m *metrics.Metric                
                switch typeM {
                case "gauge":
                        f, err := strconv.ParseFloat(value, 64)
                        if err != nil {
                                http.Error(w, "convert to gauge error", http.StatusBadRequest)
                                return
                                }
                        m = metrics.NewGaugeMetric(name, f)
        
                case "counter":
                        i, err := strconv.ParseInt(value, 10, 64)
                        if err != nil {
                                http.Error(w, "convert to counter error", http.StatusBadRequest)
                                return
                        }
                        m = metrics.NewCounterMetric(name, i)
                default:
                        http.Error(w, "type error", http.StatusNotImplemented)
                        return
                }
                
                s.AddMetric(*m)

                w.WriteHeader(http.StatusOK)
                w.Header().Set("Content-Type", "text/plain")
                fmt.Fprintln(w, m.String())
        }
}
