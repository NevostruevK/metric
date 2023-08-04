// For generation use:
// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/metric.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: proto/metric.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Metrics_GetMetric_FullMethodName       = "/metric.Metrics/GetMetric"
	Metrics_AddMetric_FullMethodName       = "/metric.Metrics/AddMetric"
	Metrics_GetAllMetrics_FullMethodName   = "/metric.Metrics/GetAllMetrics"
	Metrics_AddBatchMetrics_FullMethodName = "/metric.Metrics/AddBatchMetrics"
)

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsClient interface {
	GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
	AddMetric(ctx context.Context, in *AddMetricRequest, opts ...grpc.CallOption) (*AddMetricResponse, error)
	GetAllMetrics(ctx context.Context, in *GetAllMetricsRequest, opts ...grpc.CallOption) (*GetAllMetricsResponse, error)
	AddBatchMetrics(ctx context.Context, in *AddBatchMetricsRequest, opts ...grpc.CallOption) (*AddBatchMetricsResponse, error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, Metrics_GetMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) AddMetric(ctx context.Context, in *AddMetricRequest, opts ...grpc.CallOption) (*AddMetricResponse, error) {
	out := new(AddMetricResponse)
	err := c.cc.Invoke(ctx, Metrics_AddMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) GetAllMetrics(ctx context.Context, in *GetAllMetricsRequest, opts ...grpc.CallOption) (*GetAllMetricsResponse, error) {
	out := new(GetAllMetricsResponse)
	err := c.cc.Invoke(ctx, Metrics_GetAllMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) AddBatchMetrics(ctx context.Context, in *AddBatchMetricsRequest, opts ...grpc.CallOption) (*AddBatchMetricsResponse, error) {
	out := new(AddBatchMetricsResponse)
	err := c.cc.Invoke(ctx, Metrics_AddBatchMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility
type MetricsServer interface {
	GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
	AddMetric(context.Context, *AddMetricRequest) (*AddMetricResponse, error)
	GetAllMetrics(context.Context, *GetAllMetricsRequest) (*GetAllMetricsResponse, error)
	AddBatchMetrics(context.Context, *AddBatchMetricsRequest) (*AddBatchMetricsResponse, error)
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsServer struct {
}

func (UnimplementedMetricsServer) GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedMetricsServer) AddMetric(context.Context, *AddMetricRequest) (*AddMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMetric not implemented")
}
func (UnimplementedMetricsServer) GetAllMetrics(context.Context, *GetAllMetricsRequest) (*GetAllMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMetrics not implemented")
}
func (UnimplementedMetricsServer) AddBatchMetrics(context.Context, *AddBatchMetricsRequest) (*AddBatchMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBatchMetrics not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_GetMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).GetMetric(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_AddMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).AddMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_AddMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).AddMetric(ctx, req.(*AddMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_GetAllMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).GetAllMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_GetAllMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).GetAllMetrics(ctx, req.(*GetAllMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_AddBatchMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBatchMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).AddBatchMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_AddBatchMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).AddBatchMetrics(ctx, req.(*AddBatchMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metric.Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMetric",
			Handler:    _Metrics_GetMetric_Handler,
		},
		{
			MethodName: "AddMetric",
			Handler:    _Metrics_AddMetric_Handler,
		},
		{
			MethodName: "GetAllMetrics",
			Handler:    _Metrics_GetAllMetrics_Handler,
		},
		{
			MethodName: "AddBatchMetrics",
			Handler:    _Metrics_AddBatchMetrics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/metric.proto",
}