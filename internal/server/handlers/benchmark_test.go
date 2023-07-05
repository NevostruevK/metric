package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/go-chi/chi/v5"
)

const (
	initialBatchMetricCapacity = 200
	hashKey                    = "secretKeyForBenchmarking"
)

func Request(ts *httptest.Server, method, path string, data []byte) (int, []byte) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = resp.Body.Close()
	}()

	return resp.StatusCode, body
}

func prepareData(sM []metrics.Metrics) []byte {
	for i, m := range sM {
		if err := sM[i].SetHash(hashKey); err != nil {
			panic(m)
		}
	}
	if len(sM) == 1 {
		data, err := json.Marshal(sM[0])
		if err != nil {
			panic(err)
		}
		return data
	}
	data, err := json.Marshal(sM)
	if err != nil {
		panic(err)
	}
	return data
}

func BenchmarkRouter(b *testing.B) {
	s := storage.NewMemStorage(false, false, "")
	r := chi.NewRouter()

	r.Post("/updates/", handlers.AddBatchMetricJSONHandler(s, hashKey, initialBatchMetricCapacity))
	r.Post("/update/", handlers.AddMetricJSONHandler(s, hashKey))
	r.Post("/value/", handlers.GetMetricJSONHandler(s, hashKey))

	ts := httptest.NewServer(r)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		b.StopTimer()

		sM, _ := metrics.GetAdvanced()
		sM = append(sM, metrics.Get()...)
		data := prepareData(sM)
		mGauge := metrics.NewJSONGaugeMetric("TestGaugeMetric", 123.45)
		mCounter := metrics.NewJSONCounterMetric("PollCount", 1)
		dataGauge := prepareData([]metrics.Metrics{mGauge})
		dataCounter := prepareData([]metrics.Metrics{mCounter})

		b.StartTimer()

		if code, _ := Request(ts, "POST", "/updates/", data); code != http.StatusOK {
			panic("/updates/")
		}
		if code, _ := Request(ts, "POST", "/update/", dataGauge); code != http.StatusOK {
			panic("/update/")
		}
		if code, _ := Request(ts, "POST", "/update/", dataCounter); code != http.StatusOK {
			panic("/update/")
		}
		if code, _ := Request(ts, "POST", "/value/", dataGauge); code != http.StatusOK {
			panic("/value/")
		}
		if code, _ := Request(ts, "POST", "/value/", dataCounter); code != http.StatusOK {
			panic("/value/")
		}
	}
}
