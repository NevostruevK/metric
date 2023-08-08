package server

import (
	"context"

	"github.com/NevostruevK/metric/internal/storage"
	pb "github.com/NevostruevK/metric/proto"
	"google.golang.org/grpc/codes"
)

func (s *MetricsServer) AddBatchMetrics(ctx context.Context, in *pb.AddBatchMetricsRequest) (*pb.AddBatchMetricsResponse, error) {
	op := "gRPC AddBatchMetrics :"

	sm, err := ToSliceMetrics(in.Metrics)
	if err != nil {
		return nil, err
	}

	if err = s.checkHash(sm); err != nil {
		return nil, err
	}

	sm, err = storage.PrepareMetricsForStorage(sm)
	if err != nil {
		return nil, s.codeError(codes.Internal, err, op, "PrepareMetricsForStorage: ")
	}

	if err = s.St.AddGroupOfMetrics(ctx, sm); err != nil {
		return nil, s.codeError(codes.Internal, err, op, "AddGroupOfMetrics: ")
	}
	return &pb.AddBatchMetricsResponse{}, nil
}
