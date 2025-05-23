// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: gateway.proto

package service

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

// GatewayClient is the client API for Gateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GatewayClient interface {
	DelConn(ctx context.Context, in *GatewayRequest, opts ...grpc.CallOption) (*GatewayResponse, error)
	Push(ctx context.Context, in *GatewayRequest, opts ...grpc.CallOption) (*GatewayResponse, error)
}

type gatewayClient struct {
	cc grpc.ClientConnInterface
}

func NewGatewayClient(cc grpc.ClientConnInterface) GatewayClient {
	return &gatewayClient{cc}
}

func (c *gatewayClient) DelConn(ctx context.Context, in *GatewayRequest, opts ...grpc.CallOption) (*GatewayResponse, error) {
	out := new(GatewayResponse)
	err := c.cc.Invoke(ctx, "/service.Gateway/DelConn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayClient) Push(ctx context.Context, in *GatewayRequest, opts ...grpc.CallOption) (*GatewayResponse, error) {
	out := new(GatewayResponse)
	err := c.cc.Invoke(ctx, "/service.Gateway/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatewayServer is the server API for Gateway service.
// All implementations should embed UnimplementedGatewayServer
// for forward compatibility
type GatewayServer interface {
	DelConn(context.Context, *GatewayRequest) (*GatewayResponse, error)
	Push(context.Context, *GatewayRequest) (*GatewayResponse, error)
}

// UnimplementedGatewayServer should be embedded to have forward compatible implementations.
type UnimplementedGatewayServer struct {
}

func (UnimplementedGatewayServer) DelConn(context.Context, *GatewayRequest) (*GatewayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelConn not implemented")
}
func (UnimplementedGatewayServer) Push(context.Context, *GatewayRequest) (*GatewayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}

// UnsafeGatewayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GatewayServer will
// result in compilation errors.
type UnsafeGatewayServer interface {
	mustEmbedUnimplementedGatewayServer()
}

func RegisterGatewayServer(s grpc.ServiceRegistrar, srv GatewayServer) {
	s.RegisterService(&Gateway_ServiceDesc, srv)
}

func _Gateway_DelConn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GatewayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).DelConn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Gateway/DelConn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).DelConn(ctx, req.(*GatewayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gateway_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GatewayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Gateway/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).Push(ctx, req.(*GatewayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Gateway_ServiceDesc is the grpc.ServiceDesc for Gateway service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gateway_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Gateway",
	HandlerType: (*GatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DelConn",
			Handler:    _Gateway_DelConn_Handler,
		},
		{
			MethodName: "Push",
			Handler:    _Gateway_Push_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gateway.proto",
}
