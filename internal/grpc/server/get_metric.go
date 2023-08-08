package server

import (
	"context"

	pb "github.com/NevostruevK/metric/proto"
	"google.golang.org/grpc/codes"
)

func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	op := "gRPC Get Metric"

	mIn, err := ToMetric(in.Metric)
	if err != nil {
		return nil, err
	}
	rt, err := s.St.GetMetric(ctx, mIn.MType, mIn.ID)
	if err != nil {
		return nil, s.codeError(codes.Internal, err, op)
	}
	mOut := rt.ConvertToMetrics()

	if err = mOut.SetHash(s.Cfg.HashKey); err != nil {
		return nil, s.codeError(codes.Internal, err, op)
	}
	return &pb.GetMetricResponse{Metric: mOut.ToProto()}, nil
}
