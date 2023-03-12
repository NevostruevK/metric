package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func TestRouter(t *testing.T) {
	s := storage.NewMemStorage()
	r := chi.NewRouter()

	r.Post("/update/{typeM}/{name}/{value}", AddMetricHandler(s))
	r.Get("/value/{typeM}/{name}", GetMetricHandler(s))
	r.Get("/", GetAllMetricsHandler(s))

	ts := httptest.NewServer(r)
	defer ts.Close()

	fmt.Println("simple ok  POST update gauge")
	statusCode, body := testRequest(t, ts, "POST", "/update/gauge/testGauge/0.1234567")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "gauge/testGauge/0.123\n", body)

	fmt.Println("simple ok  POST update counter")
	statusCode, body = testRequest(t, ts, "POST", "/update/counter/testCounter/123456789")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "counter/testCounter/123456789\n", body)

	fmt.Println("simple err POST update counter with the wrong type")
	statusCode, body = testRequest(t, ts, "POST", "/update/int/testCounter/123456789")
	assert.Equal(t, http.StatusNotImplemented, statusCode)
	assert.Equal(t, "Type int is not implemented\n", body)

	fmt.Println("simple err POST param is missed")
	statusCode, body = testRequest(t, ts, "POST", "/update/counter//123456789")
	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "param is missed\n", body)

	fmt.Println("simple err POST convert with error")
	statusCode, body = testRequest(t, ts, "POST", "/update/counter/testCounter/one")
	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "convert to counter value one with an error\n", body)

	fmt.Println("simple ok  GET counter value")
	statusCode, body = testRequest(t, ts, "GET", "/value/counter/testCounter")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "123456789\n", body)

	fmt.Println("simple ok  GET gauge value")
	statusCode, body = testRequest(t, ts, "GET", "/value/gauge/testGauge")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "0.123\n", body)

	fmt.Println("simple err GET param is missed")
	statusCode, body = testRequest(t, ts, "GET", "/value//testCounter")
	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "param(s) is missed\n", body)

	fmt.Println("simple err GET value with a not implemented type")
	statusCode, body = testRequest(t, ts, "GET", "/value/int/testCounter")
	assert.Equal(t, http.StatusNotImplemented, statusCode)
	assert.Equal(t, "Type int is not implemented\n", body)

	fmt.Println("simple err GET value not found")
	statusCode, body = testRequest(t, ts, "GET", "/value/counter/unknownName")
	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, "Type counter, Name unknownName not found\n", body)
}
