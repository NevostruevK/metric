package metrics_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/NevostruevK/metric/internal/util/sign"
	"github.com/stretchr/testify/require"
)

func createKey(t *testing.T, size int) string {
	b := make([]byte, size)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return hex.EncodeToString(b)
}

func TestMetrics_CountHash(t *testing.T) {
	key := createKey(t, 16)
	gaugeName := "gaugeMetric"
	gaugeValue := 1.23456789
	gaugeHash, err := sign.Hash(fmt.Sprintf("%s:gauge:%f", gaugeName, gaugeValue), key)
	require.NoError(t, err)
	fmt.Println(gaugeHash)
	counterName := "counterMetric"
	var counterValue int64 = 123456789
	counterHash, err := sign.Hash(fmt.Sprintf("%s:counter:%d", counterName, counterValue), key)
	require.NoError(t, err)
	fmt.Println(counterHash)

	tests := []struct {
		name     string
		m        metrics.Metrics
		key      string
		wantHash string
		wantErr  bool
	}{
		{
			name:     "gauge normal test",
			m:        metrics.Metrics{gaugeName, metrics.Gauge, nil, &gaugeValue, ""},
			key:      key,
			wantHash: gaugeHash,
			wantErr:  false,
		},
		{
			name:     "countr normal test",
			m:        metrics.Metrics{counterName, metrics.Counter, &counterValue, nil, ""},
			key:      key,
			wantHash: counterHash,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, err := tt.m.CountHash(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Metrics.CountHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHash != tt.wantHash {
				t.Errorf("Metrics.CountHash() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}
