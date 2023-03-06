package client

import (
	"fmt"
	"net/http"
)
func URLHandler(w http.ResponseWriter, r *http.Request) {
    // извлекаем фрагмент query= из URL запроса search?query=something
    q := r.URL
    fmt.Println("request URL",q)
	ct:= r.Header.Get("Content-Type")
    fmt.Println("request Header",ct)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Server response")
}
/*
func TestSendMetric(t *testing.T) {
	fmt.Println("Start TestSendMetric")
	http.HandleFunc("/", URLHandler)
	fmt.Println("Run server")
    go log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("After Run server")
	
	tests := []struct {
		name string
		sM []metrics.Metric
	}{
		{
			name: "simple ok",
			sM: {metrics.Float64ToGauge(1234).NewMetric("GaugeMetric")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMetric(tt.sM)
		})
	}
}
*/