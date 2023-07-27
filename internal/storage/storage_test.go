package storage_test

import (
	"context"
	"sort"
	"testing"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddMetric(t *testing.T) {
	s := storage.NewMemStorage(false, false, "")
	var (
		nameCounter     = "testCounter"
		nameGauge       = "testGauge"
		valueCounter    = int64(1234567)
		valueCounterAdd = int64(8765432)
		valueGauge      = float64(1.234567)
		mGauge          = metrics.NewGaugeMetric(nameGauge, valueGauge)
		mCounter        = metrics.NewCounterMetric(nameCounter, valueCounter)
		mCounterAdd     = metrics.NewCounterMetric(nameCounter, valueCounterAdd)
	)
	err := s.AddMetric(context.Background(), mGauge)
	require.NoError(t, err)
	err = s.AddMetric(context.Background(), mCounter)
	require.NoError(t, err)

	assert.Equal(t, 1, len(s.Float))
	assert.Equal(t, 1, len(s.Int))
	assert.Equal(t, valueGauge, s.Float[nameGauge])
	assert.Equal(t, valueCounter, s.Int[nameCounter])

	err = s.AddMetric(context.Background(), mCounterAdd)
	require.NoError(t, err)
	assert.Equal(t, 1, len(s.Int))
	assert.Equal(t, valueCounter+valueCounterAdd, s.Int[nameCounter])
}

func TestAddGroupOfMetric(t *testing.T) {
	s := storage.NewMemStorage(false, false, "")
	sIn, err := metrics.GetAdvanced()
	require.NoError(t, err)
	sIn = append(sIn, metrics.Get()...)

	err = s.AddGroupOfMetrics(context.Background(), sIn)
	require.NoError(t, err)
	sOut, err := s.GetAllMetrics(context.Background())
	require.NoError(t, err)
	sort.Slice(sIn, func(i, j int) bool {
		return sIn[i].Name() < sIn[j].Name()
	})
	sort.Slice(sOut, func(i, j int) bool {
		return sOut[i].Name() < sOut[j].Name()
	})
	assert.Equal(t, sIn, sOut)
}
