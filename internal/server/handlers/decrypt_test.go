package handlers_test

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/crypt"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrypt(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)
	cr := &crypt.Crypt{
		PublicKey: &privateKey.PublicKey,
		Nonce:     *new([12]byte),
	}
	dcr := &crypt.Decrypt{
		PrivateKey: privateKey,
		Nonce:      *new([12]byte),
	}

	s := storage.NewMemStorage(false, false, "")
	r := chi.NewRouter()
	handlers.Logger = logger.NewLogger(`server: `, log.LstdFlags)

	r.Post("/update/", handlers.AddMetricJSONHandler(s, hashKey))

	t.Run("ok : encrypt/decrypt data", func(t *testing.T) {

		m := metrics.NewJSONGaugeMetric("testGaugeMetric", 123.4567)

		data, err := client.PrepareDataForMetric(m, hashKey, cr)
		require.NoError(t, err)

		handler := handlers.DecryptHanlder(r, dcr)
		ts := httptest.NewServer(handler)
		defer ts.Close()

		code, _ := Request(ts, "POST", "/update/", data, nil)
		assert.Equal(t, http.StatusOK, code)
		ts.Close()
	})
	t.Run("error : no command to decrypt data", func(t *testing.T) {

		m := metrics.NewJSONGaugeMetric("testGaugeMetric", 123.4567)

		data, err := client.PrepareDataForMetric(m, hashKey, cr)
		require.NoError(t, err)

		handler := handlers.DecryptHanlder(r, nil)
		ts := httptest.NewServer(handler)
		defer ts.Close()

		code, _ := Request(ts, "POST", "/update/", data, nil)
		assert.Equal(t, http.StatusBadRequest, code)
		ts.Close()
	})
	t.Run("error : no command to encrypt data", func(t *testing.T) {

		m := metrics.NewJSONGaugeMetric("testGaugeMetric", 123.4567)

		data, err := client.PrepareDataForMetric(m, hashKey, nil)
		require.NoError(t, err)

		handler := handlers.DecryptHanlder(r, dcr)
		ts := httptest.NewServer(handler)
		defer ts.Close()

		code, _ := Request(ts, "POST", "/update/", data, nil)
		assert.Equal(t, http.StatusBadRequest, code)
		ts.Close()
	})
	t.Run("ok : without encrypting and decrypting data", func(t *testing.T) {
		m := metrics.NewJSONGaugeMetric("testGaugeMetric", 123.4567)

		data, err := client.PrepareDataForMetric(m, hashKey, nil)
		require.NoError(t, err)

		handler := handlers.DecryptHanlder(r, nil)
		ts := httptest.NewServer(handler)
		defer ts.Close()

		code, _ := Request(ts, "POST", "/update/", data, nil)
		assert.Equal(t, http.StatusOK, code)
		ts.Close()
	})

}
