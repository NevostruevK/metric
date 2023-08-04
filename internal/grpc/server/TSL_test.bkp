package server_test

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/NevostruevK/metric/proto"

	grpcserver "github.com/NevostruevK/metric/internal/grpc/server"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TLS(t *testing.T) {
	cfg := &commands.Config{}
	cfg.SetOption(commands.WithAddress("127.0.0.1:8084"))
	t.Run("normal work without TLS", func(t *testing.T) {
		cfg := &commands.Config{}
		cfg.SetOption(commands.WithAddress("127.0.0.1:8084"))
		conn, client, err := startClient(cfg)
		require.NoError(t, err)
		s, err := startServer(cfg)
		require.NoError(t, err)
		go s.ListenAndServe()

		mIn := metrics.NewJSONGaugeMetric("test_gauge", 123.45)
		resp, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		mOut, _ := grpcserver.ToMetric(resp.Metric)
		assert.Equal(t, mIn, *mOut)
		defer conn.Close()
		s.Srv.GracefulStop()
	})

	t.Run(`error server use TLS, agent dosn't`, func(t *testing.T) {
		cfg := &commands.Config{}
		cfg.SetOption(commands.WithAddress("127.0.0.1:8085"))
		conn, client, err := startClient(cfg)
		require.NoError(t, err)
		cfg.SetOption(commands.WithCertificate("./data/server.crt"))
		cfg.SetOption(commands.WithCryptoKey("./data/private.key"))
		s, err := startServer(cfg)
		require.NoError(t, err)
		go s.ListenAndServe()

		mIn := metrics.NewJSONGaugeMetric("test_gauge", 123.45)
		_, err = client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		fmt.Println(err)
		require.Error(t, err)
		defer conn.Close()
		s.Srv.GracefulStop()
	})

	t.Run(`error client use TLS, server dosn't`, func(t *testing.T) {
		cfg := &commands.Config{}
		cfg.SetOption(commands.WithAddress("127.0.0.1:8086"))
		cfg.SetOption(commands.WithCertificate("./data/server.crt"))
		cfg.SetOption(commands.WithCryptoKey(""))
		conn, client, err := startClient(cfg)
		require.NoError(t, err)
		s, err := startServer(cfg)
		require.NoError(t, err)
		go s.ListenAndServe()

		mIn := metrics.NewJSONGaugeMetric("test_gauge", 123.45)
		_, err = client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		fmt.Println(err)
		require.Error(t, err)
		defer conn.Close()
		s.Srv.GracefulStop()
	})

	t.Run("normal work with TLS", func(t *testing.T) {
		cfg := &commands.Config{}
		cfg.SetOption(commands.WithAddress("127.0.0.1:8086"))
		cfg.SetOption(commands.WithCertificate("./data/server.crt"))
		cfg.SetOption(commands.WithCryptoKey("./data/private.key"))
		conn, client, err := startClient(cfg)
		require.NoError(t, err)
		s, err := startServer(cfg)
		require.NoError(t, err)
		go s.ListenAndServe()

		mIn := metrics.NewJSONGaugeMetric("test_gauge", 123.45)
		resp, err := client.AddMetric(
			context.Background(),
			&pb.AddMetricRequest{Metric: mIn.ToProto()},
		)
		require.NoError(t, err)
		mOut, _ := grpcserver.ToMetric(resp.Metric)
		assert.Equal(t, mIn, *mOut)
		defer conn.Close()
		s.Srv.GracefulStop()
	})
}
