package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/go-chi/chi/v5"
)

var Logger = &log.Logger{}

// GetPingHandler обработчик запроса /ping.
// Проверка подключения к базе.
func GetPingHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if err := s.Ping(); err != nil {
			msg := fmt.Sprintf("Can't Ping database %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ping database ok")
	}
}

// GetAllMetricsHandler обработчик запроса /.
// Выдача всех метрик.
func GetAllMetricsHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm, err := s.GetAllMetrics(context.Background())
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			msg := fmt.Sprintf("ERROR : GetAllMetricsHandler:GetAllMetrics() returned the error %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		var prefix = `<html class="h-100">
    			<head>
    				<title></title>
    			</head>
    			<body class="d-flex flex-column h-100">
					<table class="table">
						<th>type</th>
						<th>name</th>
						<th>value</th>`
		var suffix = `</table>
				</body>
			</html>`

		w.Write([]byte(prefix))
		for _, v := range sm {
			w.Write([]byte("<tr><td>" + v.Name() + "</td>"))
			w.Write([]byte("<td>" + v.Type() + "</td>"))
			w.Write([]byte("<td>" + v.StringValue() + "</td></tr>"))
		}
		w.Write([]byte(suffix))
	}
}

// GetMetricHandler обработчик запроса /value/{typeM}/{name}.
// Выдача метрики name типа typeM.
func GetMetricHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeM := chi.URLParam(r, "typeM")
		name := chi.URLParam(r, "name")
		if typeM == "" || name == "" {
			msg := fmt.Sprintln("ERROR : GetMetricHandler param(s) is missed ")
			Logger.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		if isValidType := metrics.IsMetricType(typeM); !isValidType {
			msg := fmt.Sprintf("type %s is not implemented", typeM)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusNotImplemented)
			return
		}

		rt, err := s.GetMetric(context.Background(), typeM, name)
		if err != nil {
			msg := fmt.Sprintf("ERROR : GetMetricHandler:GetMetric() returned the error %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, rt.StringValue())
	}
}

// AddMetricHandler обработчик запроса /update/{typeM}/{name}/{value}.
// Прием метрики name типа typeM со значением value.
func AddMetricHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeM := chi.URLParam(r, "typeM")
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")
		if typeM == "" || name == "" || value == "" {
			msg := fmt.Sprintln("ERROR : AddMetricHandler param(s) is missed")
			Logger.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		if isValidType := metrics.IsMetricType(typeM); !isValidType {
			msg := fmt.Sprintf("type %s is not implemented", typeM)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusNotImplemented)
			return
		}

		m, err := metrics.NewValueMetric(name, typeM, value)
		if err != nil {
			msg := fmt.Sprintf("ERROR : AddMetricHandler:metrics.NewValueMetric() returned the error %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return

		}

		if err = s.AddMetric(context.Background(), m); err != nil {
			msg := fmt.Sprintf("ERROR : AddMetricHandler:metrics.AddMetric(m) returned the error %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, m.StringWithSlash())
	}
}
