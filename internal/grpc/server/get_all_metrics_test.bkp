package server_test

import (
	"context"
	"testing"

	grpcserver "github.com/NevostruevK/metric/internal/grpc/server"
	pb "github.com/NevostruevK/metric/proto"

	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsServer_GetAllMetrics(t *testing.T) {
	cfg := &commands.Config{}
	cfg.SetOption(commands.WithAddress("127.0.0.1:8083"))
	conn, client, err := startClient(cfg)
	require.NoError(t, err)
	defer conn.Close()
	s, err := startServer(cfg)
	require.NoError(t, err)
	go s.ListenAndServe()

	t.Run("Get Empty slice of metrics", func(t *testing.T) {

		resp, err := client.GetAllMetrics(
			context.Background(),
			&pb.GetAllMetricsRequest{},
		)
		require.NoError(t, err)
		smOut, err := grpcserver.ToSliceMetrics(resp.Metrics)
		require.NoError(t, err)
		assert.Empty(t, smOut)
	})

	t.Run("Get all metrics", func(t *testing.T) {
		smIn := []metrics.Metrics{
			metrics.NewJSONGaugeMetric("test_gauge", 123.45),
			metrics.NewJSONGaugeMetric("test_gauge1", 789.45),
			metrics.NewJSONCounterMetric("test_counter", 12345),
			metrics.NewJSONCounterMetric("test_counter1", 67890),
		}
		_, err := client.AddBatchMetrics(
			context.Background(),
			&pb.AddBatchMetricsRequest{Metrics: metrics.ToProto(smIn)},
		)
		require.NoError(t, err)

		resp, err := client.GetAllMetrics(
			context.Background(),
			&pb.GetAllMetricsRequest{},
		)
		require.NoError(t, err)
		smOut, err := grpcserver.ToSliceMetrics(resp.Metrics)
		require.NoError(t, err)
		assert.ElementsMatch(t, smIn, smOut)
	})
}
