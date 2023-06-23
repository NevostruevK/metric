package handlers_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
)

func ExampleGetMetricHandler() {
	s := storage.NewMemStorage(false, false, "")
	r := chi.NewRouter()

	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", handlers.GetMetricHandler(s))

	ts := httptest.NewServer(r)
	defer ts.Close()

	Request(ts, "POST", "/update/gauge/exampleGauge/1.2345", nil)
	Request(ts, "POST", "/update/counter/exampleCounter/12345", nil)

	code, body := Request(ts, "GET", "/value/gauge/exampleGauge", nil)
	fmt.Printf("%d : %s", code, string(body))

	code, body = Request(ts, "GET", "/value/counter/exampleCounter", nil)
	fmt.Printf("%d : %s", code, string(body))

	// Output:
	// 200 : 1.234
	// 200 : 12345
}
