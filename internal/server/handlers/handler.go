package handlers

import (
	"compress/gzip"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/go-chi/chi/v5"
)

type gzipWriter struct {
    http.ResponseWriter
    Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
    // w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
    return w.Writer.Write(b)
} 

func CompressHandle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // проверяем, что клиент поддерживает gzip-сжатие
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            // если gzip не поддерживается, передаём управление
            // дальше без изменений
            next.ServeHTTP(w, r)
            return
        }

        // создаём gzip.Writer поверх текущего w
        gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
        if err != nil {
            io.WriteString(w, err.Error())
            return
        }
        defer gz.Close()

        w.Header().Set("Content-Encoding", "gzip")
        // передаём обработчику страницы переменную типа gzipWriter для вывода данных
        next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
    })
}

func DecompressHanlder(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			io.WriteString(w, err.Error())
			return 
		}
		defer gz.Close()	
		newReq := r
		newReq.Body = gz
		next.ServeHTTP(w, newReq)		
	})
}

func GetAllMetricsHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm := s.GetAllMetrics()
		path := filepath.Join(".", "..", "..", "internal", "files", "html", "getAllMetrics.html")
		//создаем html-шаблон
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//исполняем именованный шаблон "Metric", передавая туда массив со списком метрик
		err = tmpl.ExecuteTemplate(w, "Metric", sm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

		rt, err := s.GetMetric(typeM, name)
		if err != nil {
			http.Error(w, "Type "+typeM+", Name "+name+" not found", http.StatusNotFound)
			fmt.Println("Type ",typeM,", Name ",name," not found")
			return
		}
		m, ok := rt.(*metrics.Metric)
		if !ok{
			http.Error(w, "Type "+typeM+", Name "+name+" is not a metric type", http.StatusNotFound)
			fmt.Println("Type ",typeM,", Name ",name," is not a metric type")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, m.StringValue())
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
			fmt.Println("AddMetricHandler : Type ",typeM,", Name ",name, " value ",value, " with an error")
			return

		}

		s.AddMetric(m)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, m.String())
	}
}
