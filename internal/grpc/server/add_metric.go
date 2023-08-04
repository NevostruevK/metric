package server

import (
	"context"

	"github.com/NevostruevK/metric/internal/util/metrics"
	pb "github.com/NevostruevK/metric/proto"
	"google.golang.org/grpc/codes"
)

func (s *MetricsServer) AddMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	op := "gRPC Add Metric"

	m, err := ToMetric(in.Metric)
	if err != nil {
		return nil, err
	}
	if err = s.checkHash([]metrics.Metrics{*m}); err != nil {
		return nil, err
	}
	if err := s.St.AddMetric(ctx, m); err != nil {
		return nil, s.codeError(codes.Internal, err, op)
	}
	if err = m.SetHash(s.Cfg.HashKey); err != nil {
		return nil, s.codeError(codes.Internal, err, op)
	}
	return &pb.AddMetricResponse{Metric: m.ToProto()}, nil
}
