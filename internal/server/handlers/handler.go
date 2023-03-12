package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/go-chi/chi/v5"
)

func GetAllMetricsHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm := s.GetAllMetrics()
		path := filepath.Join(".", "..", "..", "internal", "files", "html", "getAllMetrics.html")
		//создаем html-шаблон
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		//исполняем именованный шаблон "Metric", передавая туда массив со списком метрик
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
		if typeM == "" || name == "" {
			http.Error(w, "param(s) is missed", http.StatusBadRequest)
			return
		}
		if isValidType := metrics.IsMetricType(typeM); !isValidType {
			http.Error(w, "Type "+typeM+" is not implemented", http.StatusNotImplemented)
			return
		}

		m, err := s.GetMetric(typeM, name)
		if err != nil {
			http.Error(w, "Type "+typeM+", Name "+name+" not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, m.Value())
	}
}

func AddMetricHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeM := chi.URLParam(r, "typeM")
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")
		if typeM == "" || name == "" || value == "" {
			http.Error(w, "param is missed", http.StatusBadRequest)
			return
		}

		if isValidType := metrics.IsMetricType(typeM); !isValidType {
			http.Error(w, "Type "+typeM+" is not implemented", http.StatusNotImplemented)
			return
		}

		m, err := metrics.NewValueMetric(name, typeM, value)
		if err != nil {
			http.Error(w, "convert to "+typeM+" value "+value+" with an error", http.StatusBadRequest)
			return

		}

		s.AddMetric(*m)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, m.String())
	}
}
