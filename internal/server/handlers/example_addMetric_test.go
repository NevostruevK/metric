package handlers_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
)

func ExampleAddMetricHandler() {
	s := storage.NewMemStorage(false, false, "")
	r := chi.NewRouter()

	r.Post("/update/{typeM}/{name}/{value}", handlers.AddMetricHandler(s))

	ts := httptest.NewServer(r)
	defer ts.Close()

	code, body := Request(ts, "POST", "/update/gauge/exampleGauge/1.2345", nil, nil)
	fmt.Printf("%d : %s", code, string(body))

	code, body = Request(ts, "POST", "/update/counter/exampleCounter/12345", nil, nil)
	fmt.Printf("%d : %s", code, string(body))

	// Output:
	// 200 : gauge/exampleGauge/1.234
	// 200 : counter/exampleCounter/12345
}
