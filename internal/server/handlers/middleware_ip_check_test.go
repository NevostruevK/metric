package handlers_test

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPCheckHandler(t *testing.T) {
	const cidr = "192.168.1.0/24"
	_, ipNet, err := net.ParseCIDR(cidr)
	require.NoError(t, err)

	m := metrics.NewJSONGaugeMetric("testGaugeMetric", 123.4567)
	data, err := client.PrepareDataForMetric(m, hashKey, nil)
	require.NoError(t, err)

	s := storage.NewMemStorage(false, false, "")

	r := chi.NewRouter()
	handlers.Logger = logger.NewLogger(`server: `, log.LstdFlags)
	handler := handlers.IPCheckHandler(r, ipNet)
	r.Post("/update/", handlers.AddMetricJSONHandler(s, hashKey))

	t.Run("ok : subnet contain ip", func(t *testing.T) {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		code, _ := Request(ts, "POST", "/update/", data, map[string]string{"X-Real-IP": "192.168.1.128"})
		assert.Equal(t, http.StatusOK, code)
		ts.Close()
	})
	t.Run("Forbidden : subnet doesn't contain ip", func(t *testing.T) {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		code, _ := Request(ts, "POST", "/update/", data, map[string]string{"X-Real-IP": "192.168.2.0"})
		assert.Equal(t, http.StatusForbidden, code)
		ts.Close()
	})
	t.Run(`Bad request : there is no "X-Real-IP" header`, func(t *testing.T) {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		code, _ := Request(ts, "POST", "/update/", data, nil)
		assert.Equal(t, http.StatusBadRequest, code)
		ts.Close()
	})

}
