// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: otter/v1/otter.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Otter_Status_FullMethodName = "/otter.v1.Otter/Status"
)

// OtterClient is the client API for Otter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OtterClient interface {
	// Implements a client-side heartbeat that can also be used by monitoring tools.
	Status(ctx context.Context, in *HealthCheck, opts ...grpc.CallOption) (*ServiceState, error)
}

type otterClient struct {
	cc grpc.ClientConnInterface
}

func NewOtterClient(cc grpc.ClientConnInterface) OtterClient {
	return &otterClient{cc}
}

func (c *otterClient) Status(ctx context.Context, in *HealthCheck, opts ...grpc.CallOption) (*ServiceState, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ServiceState)
	err := c.cc.Invoke(ctx, Otter_Status_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OtterServer is the server API for Otter service.
// All implementations must embed UnimplementedOtterServer
// for forward compatibility
type OtterServer interface {
	// Implements a client-side heartbeat that can also be used by monitoring tools.
	Status(context.Context, *HealthCheck) (*ServiceState, error)
	mustEmbedUnimplementedOtterServer()
}

// UnimplementedOtterServer must be embedded to have forward compatible implementations.
type UnimplementedOtterServer struct {
}

func (UnimplementedOtterServer) Status(context.Context, *HealthCheck) (*ServiceState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedOtterServer) mustEmbedUnimplementedOtterServer() {}

// UnsafeOtterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OtterServer will
// result in compilation errors.
type UnsafeOtterServer interface {
	mustEmbedUnimplementedOtterServer()
}

func RegisterOtterServer(s grpc.ServiceRegistrar, srv OtterServer) {
	s.RegisterService(&Otter_ServiceDesc, srv)
}

func _Otter_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheck)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OtterServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Otter_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OtterServer).Status(ctx, req.(*HealthCheck))
	}
	return interceptor(ctx, in, info, handler)
}

// Otter_ServiceDesc is the grpc.ServiceDesc for Otter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Otter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "otter.v1.Otter",
	HandlerType: (*OtterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _Otter_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "otter/v1/otter.proto",
}