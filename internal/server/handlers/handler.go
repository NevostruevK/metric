package handlers

import (
	"fmt"
	"net/http"

	"github.com/NevostruevK/metric/internal/db"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/go-chi/chi/v5"
)

func GetPingHandler(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if err := db.Ping(); err != nil {
			logger.LogError("GetPingHandler : ", "db.Ping", err)
			http.Error(w, fmt.Sprintf("Can't Ping database %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ping database ok")
	}
}

func GetAllMetricsHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm, err := s.GetAllMetrics()
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			logger.LogError("GetAllMetricsHandler : ", "s.GetAllMetrics()", err)
			http.Error(w, fmt.Sprintf("Can't get all metrics  %v", err), http.StatusInternalServerError)
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

func GetMetricHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeM := chi.URLParam(r, "typeM")
		name := chi.URLParam(r, "name")
		if typeM == "" || name == "" {
			logger.LogError("GetMetricHandler : ", "param(s) is missed", nil)
			http.Error(w, "param(s) is missed", http.StatusBadRequest)
			return
		}
		if isValidType := metrics.IsMetricType(typeM); !isValidType {
			msg := fmt.Sprintf("type %s is not implemented", typeM)
			logger.LogError("GetMetricHandler : ", msg, nil)
			http.Error(w, msg, http.StatusNotImplemented)
			return
		}

		rt, err := s.GetMetric(typeM, name)
		if err != nil {
			logger.LogError("GetMetricHandler : ", "s.GetMetric()", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, rt.StringValue())
	}
}

func AddMetricHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeM := chi.URLParam(r, "typeM")
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")
		if typeM == "" || name == "" || value == "" {
			logger.LogError("AddMetricHandler : ", "param(s) is missed", nil)
			http.Error(w, "param is missed", http.StatusBadRequest)
			return
		}

		if isValidType := metrics.IsMetricType(typeM); !isValidType {
			msg := fmt.Sprintf("type %s is not implemented", typeM)
			logger.LogError("AddMetricHandler : ", msg, nil)
			http.Error(w, msg, http.StatusNotImplemented)
			return
		}

		m, err := metrics.NewValueMetric(name, typeM, value)
		if err != nil {
			logger.LogError("AddMetricHandler : ", "metrics.NewValueMetric", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		if err = s.AddMetric(m); err != nil {
			logger.LogError("AddMetricHandler : ", "s.AddMetric", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, m.String())
	}
}
