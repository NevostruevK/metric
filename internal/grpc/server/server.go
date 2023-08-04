package server

import (
	// импортируем пакет со сгенерированными protobuf-файлами

	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
	pb "github.com/NevostruevK/metric/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	FullMethodAddBatchMetrics = "/metric.Metrics/AddBatchMetrics"
	FullMethodAddMetric       = "/metric.Metrics/AddMetric"
	FullMethodGetMetric       = "/metric.Metrics/GetMetric"
	FullMethodGetAllMetric    = "/metric.Metrics/GetAllMetric"
)
var (
	ErrInMetricIsNil = errors.New("field Metric is nil")
	ErrWrongHash = errors.New("wrong hash for metric")
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	Srv *grpc.Server
	St  storage.Repository
	Cfg *commands.Config
	lgr *log.Logger
}

func NewServer(s storage.Repository, cfg *commands.Config) (*MetricsServer, error) {
	opts := []grpc.ServerOption{}
	if cfg.Certificate != "" && cfg.CryptoKey != "" {
		cert, err := tls.LoadX509KeyPair(cfg.Certificate, cfg.CryptoKey)
		if err != nil {
			return nil, err
		}
		fmt.Println("Start with TSL")
		opts = append(opts, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	}
	server := MetricsServer{
		Srv: grpc.NewServer(opts...),
		St:  s,
		Cfg: cfg,
		lgr: logger.NewLogger("gRPC server:", log.LstdFlags|log.Lshortfile),
	}
	return &server, nil
}

func (s *MetricsServer) ListenAndServe() error {
	listen, err := net.Listen("tcp", s.Cfg.Address)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterMetricsServer(s.Srv, s)
	return s.Srv.Serve(listen)
}

func (s *MetricsServer) Shutdown(ctx context.Context) error {
	s.Srv.GracefulStop()
	return nil
}

func ToMetric(m *pb.Metric) (*metrics.Metrics, error) {
	if m == nil {
		return nil, status.Error(codes.InvalidArgument, ErrInMetricIsNil.Error())
	}
	var metric metrics.Metrics
	if m.Type == pb.MetricType_GAUGE {
		metric = metrics.NewJSONGaugeMetric(m.Name, m.Value)
	} else {
		metric = metrics.NewJSONCounterMetric(m.Name, m.Delta)
	}
	metric.Hash = m.Hash
	return &metric, nil
}

func ToSliceMetrics(pm []*pb.Metric) ([]metrics.Metrics, error) {
	sm := make([]metrics.Metrics, len(pm))
	for i, pm := range pm {
		m, err := ToMetric(pm)
		if err != nil {
			return nil, err
		}
		sm[i] = *m
	}
	return sm, nil
}

func (s *MetricsServer) codeError(code codes.Code, err error, msg ...string) error {
	text := strings.Join(msg, "") + fmt.Sprintf(" failed with error %v", err)
	s.lgr.Println(text)
	return status.Error(code, text)
}

func (s *MetricsServer) checkHash(sm []metrics.Metrics) error {
	if s.Cfg.HashKey == "" {
		return nil
	}
	for _, m := range sm {
		ok, err := m.CheckHash(s.Cfg.HashKey)
		if err != nil {
			return s.codeError(codes.Internal, err, m.String())
		}
		if !ok {
			return s.codeError(codes.InvalidArgument, ErrWrongHash, m.String())
		}
	}
	return nil
}

func (s *MetricsServer) setHash(sm []metrics.Metrics) error {
	if s.Cfg.HashKey == "" {
		return nil
	}
	for _, m := range sm {
		if err := m.SetHash(s.Cfg.HashKey); err != nil {
			return err
		}
	}
	return nil
}
