package handlers

import (
	"fmt"
	"net/http"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/go-chi/chi/v5"
)

func GetAllMetricsHandler(s storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm := s.GetAllMetrics()
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

		s.AddMetric(m)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, m.String())
	}
}
