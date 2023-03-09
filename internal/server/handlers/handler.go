package handlers

import (
	"fmt"
	"net/http"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

func URLHandler(w http.ResponseWriter, r *http.Request) {
        // извлекаем фрагмент query= из URL запроса search?query=something
        //    q := r.URL
        //    fmt.Println("request URL",q)
        //      ct:= r.Header.Get("Content-Type")
        //    fmt.Println("request Header",ct)
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Server response")
}

func SaveMetricHandler(s *storage.MemStorage) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {

                url := r.URL.String()
                m, err := metrics.URLToMetric(url)
                if err != nil {
                        http.Error(w, "can't parse URL", http.StatusNotFound)
                        return
                }
                s.SaveMetric(*m)

                w.WriteHeader(http.StatusOK)
                w.Header().Set("Content-Type", "text/plain")
                fmt.Fprintln(w, m.String())
        }
}
func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
        http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}