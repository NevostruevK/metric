package server_test

import (
	"context"
	"testing"

	pb "github.com/NevostruevK/metric/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpcserver "github.com/NevostruevK/metric/internal/grpc/server"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsServer_AddMetric(t *testing.T) {
	cfg := &commands.Config{}
	cfg.SetOption(commands.WithAddress("127.0.0.1:8080"))
	conn, client, err := startClient(cfg)
	require.NoError(t, err)
	defer conn.Close()
	s, err := startServer(cfg)
	require.NoError(t, err)
	go s.ListenAndServe()

	t.Run("Add gauge metric", func(t *testing.T) {
		mIn := metrics.NewJSONGaugeMetric("test_gauge", 123.45)
		resp, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		mOut, _ := grpcserver.ToMetric(resp.Metric)
		assert.Equal(t, mIn, *mOut)
	})
	t.Run("Add counter metric", func(t *testing.T) {
		mIn := metrics.NewJSONCounterMetric("test_counter", 12345)
		resp, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		mOut, _ := grpcserver.ToMetric(resp.Metric)
		assert.Equal(t, mIn, *mOut)
	})
	t.Run("Add nil metric", func(t *testing.T) {
		_, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{},
		)
		require.Error(t, err)
		assert.Equal(t, err, status.Error(codes.InvalidArgument, grpcserver.ErrInMetricIsNil.Error()))
	})
	t.Run("Add metric with wrong hash", func(t *testing.T) {
		cfg.SetOption(commands.WithHashKey("test_hash_key"))
		mIn := &metrics.Metrics{}
		mIn.NewGaugeMetric("test_gauge", 123.45)
		mIn.Hash = "wrongHash"
		_, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, e.Code(), codes.InvalidArgument)
	})
	t.Run("Add gauge metric with hash key", func(t *testing.T) {
		cfg.SetOption(commands.WithHashKey("test_hash_key"))
		mIn := &metrics.Metrics{}
		mIn.NewGaugeMetric("test_gauge", 123.45)
		err := mIn.SetHash(cfg.HashKey)
		require.NoError(t, err)
		resp, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		mOut, _ := grpcserver.ToMetric(resp.Metric)
		assert.Equal(t, *mIn, *mOut)
	})
	t.Run("Add counter metric with hash key", func(t *testing.T) {
		cfg.SetOption(commands.WithHashKey("test_hash_key"))
		mIn := &metrics.Metrics{}
		mIn.NewCounterMetric("test_counter3", 123)
		err := mIn.SetHash(cfg.HashKey)
		require.NoError(t, err)
		resp, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		mOut, _ := grpcserver.ToMetric(resp.Metric)
		assert.Equal(t, *mIn, *mOut)
	})
	t.Run("Add the same counter metric with hash key", func(t *testing.T) {
		const (
			metricsName = "test_counter1"
			value1      = 123
			value2      = 877
		)
		cfg.SetOption(commands.WithHashKey("test_hash_key"))
		m1 := metrics.NewJSONCounterMetric(metricsName, value1)
		m2 := metrics.NewJSONCounterMetric(metricsName, value2)
		var mIn = []*metrics.Metrics{&m1, &m2}

		var resp *pb.AddMetricResponse
		for _, m := range mIn {
			err := m.SetHash(cfg.HashKey)
			require.NoError(t, err)
			resp, err = client.AddMetric(
				context.Background(),
				&pb.AddMetricRequest{Metric: m.ToProto()},
			)
			require.NoError(t, err)
		}
		mOut, _ := grpcserver.ToMetric(resp.Metric)
		mWaited := &metrics.Metrics{}
		mWaited.NewCounterMetric(metricsName, value1+value2)
		err = mWaited.SetHash(cfg.HashKey)
		require.NoError(t, err)
		assert.Equal(t, *mOut, *mWaited)
	})

}
