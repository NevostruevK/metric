package handlers_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
)

func ExampleGetAllMetricsHandler() {
	s := storage.NewMemStorage(false, false, "")
	r := chi.NewRouter()

	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/", handlers.GetAllMetricsHandler(s))

	ts := httptest.NewServer(r)
	defer ts.Close()

	Request(ts, "POST", "/update/gauge/exampleGauge/1.2345", nil)
	Request(ts, "POST", "/update/counter/exampleCounter/12345", nil)

	code, _ := Request(ts, "GET", "/", nil)
	fmt.Println(code)

	// Output:
	// 200
}
