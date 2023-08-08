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

func TestMetricsServer_GetMetric(t *testing.T) {
	cfg := &commands.Config{}
	cfg.SetOption(commands.WithAddress("127.0.0.1:8081"))
	conn, client, err := startClient(cfg)
	require.NoError(t, err)
	defer conn.Close()
	s, err := startServer(cfg)
	require.NoError(t, err)
	go s.ListenAndServe()

	t.Run("Get gauge metric", func(t *testing.T) {
		mIn := metrics.NewJSONGaugeMetric("test_gauge", 123.45)
		_, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		resp, err := client.GetMetric(
			context.Background(),
			&pb.GetMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		mOut, _ := grpcserver.ToMetric(resp.Metric)
		assert.Equal(t, mIn, *mOut)
	})
	t.Run("Get counter metric", func(t *testing.T) {
		mIn := metrics.NewJSONCounterMetric("test_counter", 12345)
		_, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		resp, err := client.GetMetric(
			context.Background(),
			&pb.GetMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		mOut, _ := grpcserver.ToMetric(resp.Metric)
		assert.Equal(t, mIn, *mOut)
	})
	t.Run("Get nil metric in parameter", func(t *testing.T) {
		_, err := client.GetMetric(
			context.Background(),
			&pb.GetMetricRequest{},
		)
		require.Error(t, err)
		assert.Equal(t, err, status.Error(codes.InvalidArgument, grpcserver.ErrInMetricIsNil.Error()))
	})
	t.Run("Get not used metric", func(t *testing.T) {
		mIn := metrics.NewJSONCounterMetric("not_used_metric", 1)
		_, err := client.GetMetric(
			context.Background(),
			&pb.GetMetricRequest{Metric: mIn.ToProto()},
		)
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, e.Code(), codes.Internal)
	})
}
