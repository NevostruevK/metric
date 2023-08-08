package server

import (
	"context"

	"github.com/NevostruevK/metric/internal/util/metrics"
	pb "github.com/NevostruevK/metric/proto"
	"google.golang.org/grpc/codes"
)

func (s *MetricsServer) GetAllMetrics(ctx context.Context, in *pb.GetAllMetricsRequest) (*pb.GetAllMetricsResponse, error) {
	op := "gRPC Get All Metrics"

	sm, err := s.St.GetAllMetrics(ctx)
	if err != nil {
		return nil, s.codeError(codes.Internal, err, op)
	}

	if err = s.setHash(sm); err != nil {
		return nil, s.codeError(codes.Internal, err, op)
	}

	return &pb.GetAllMetricsResponse{Metrics: metrics.ToProto(sm)}, nil
}
