package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

//func SaveMetricHandler(s *storage.MemStorage) http.HandlerFunc {
func AddMetricHandler(s storage.Repository) http.HandlerFunc {
                return func(w http.ResponseWriter, r *http.Request) {

                url := r.URL.String()
                fmt.Println("Get URL: ",url)
//                m, err := metrics.URLToMetric(url)
                words := strings.Split(url, "/")
                if len(words) != 5 {
                        http.Error(w, "wrong slash count error", http.StatusNotFound)
                        return
                }
                var m *metrics.Metric                
                switch words[2] {
                case "gauge":
                        f, err := strconv.ParseFloat(words[4], 64)
                        if err != nil {
                                http.Error(w, "convert to gauge error", http.StatusBadRequest)
                                return
                                }
                        m = metrics.NewGaugeMetric(words[3], f)
        
                case "counter":
                        i, err := strconv.ParseInt(words[4], 10, 64)
                        if err != nil {
                                http.Error(w, "convert to counter error", http.StatusBadRequest)
                                return
                        }
                        m = metrics.NewCounterMetric(words[3], i)
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
