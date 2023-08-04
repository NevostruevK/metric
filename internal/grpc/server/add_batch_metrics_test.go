package server_test

import (
	"context"
	"testing"

	pb "github.com/NevostruevK/metric/proto"

	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsServer_AddBatchMetrics(t *testing.T) {
	cfg := &commands.Config{}
	cfg.SetOption(commands.WithAddress("127.0.0.1:8082"))
	conn, client, err := startClient(cfg)
	require.NoError(t, err)
	defer conn.Close()
	s, err := startServer(cfg)
	require.NoError(t, err)
	go s.ListenAndServe()

	t.Run("Add metrics", func(t *testing.T) {
		sm := []metrics.Metrics{
			metrics.NewJSONGaugeMetric("test_gauge", 123.45),
			metrics.NewJSONCounterMetric("test_counter", 12345),
		}
		_, err := client.AddBatchMetrics(
			context.Background(),
			&pb.AddBatchMetricsRequest{Metrics: metrics.ToProto(sm)},
		)
		assert.NoError(t, err)
	})
}
