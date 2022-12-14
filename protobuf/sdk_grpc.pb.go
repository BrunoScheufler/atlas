// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: sdk.proto

package protobuf

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

// AtlasfileClient is the client API for Atlasfile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AtlasfileClient interface {
	Eval(ctx context.Context, in *EvalRequest, opts ...grpc.CallOption) (*EvalReply, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error)
}

type atlasfileClient struct {
	cc grpc.ClientConnInterface
}

func NewAtlasfileClient(cc grpc.ClientConnInterface) AtlasfileClient {
	return &atlasfileClient{cc}
}

func (c *atlasfileClient) Eval(ctx context.Context, in *EvalRequest, opts ...grpc.CallOption) (*EvalReply, error) {
	out := new(EvalReply)
	err := c.cc.Invoke(ctx, "/sdk.Atlasfile/Eval", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *atlasfileClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error) {
	out := new(PingReply)
	err := c.cc.Invoke(ctx, "/sdk.Atlasfile/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AtlasfileServer is the server API for Atlasfile service.
// All implementations must embed UnimplementedAtlasfileServer
// for forward compatibility
type AtlasfileServer interface {
	Eval(context.Context, *EvalRequest) (*EvalReply, error)
	Ping(context.Context, *PingRequest) (*PingReply, error)
	mustEmbedUnimplementedAtlasfileServer()
}

// UnimplementedAtlasfileServer must be embedded to have forward compatible implementations.
type UnimplementedAtlasfileServer struct {
}

func (UnimplementedAtlasfileServer) Eval(context.Context, *EvalRequest) (*EvalReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Eval not implemented")
}
func (UnimplementedAtlasfileServer) Ping(context.Context, *PingRequest) (*PingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedAtlasfileServer) mustEmbedUnimplementedAtlasfileServer() {}

// UnsafeAtlasfileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AtlasfileServer will
// result in compilation errors.
type UnsafeAtlasfileServer interface {
	mustEmbedUnimplementedAtlasfileServer()
}

func RegisterAtlasfileServer(s grpc.ServiceRegistrar, srv AtlasfileServer) {
	s.RegisterService(&Atlasfile_ServiceDesc, srv)
}

func _Atlasfile_Eval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EvalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtlasfileServer).Eval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.Atlasfile/Eval",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtlasfileServer).Eval(ctx, req.(*EvalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Atlasfile_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtlasfileServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.Atlasfile/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtlasfileServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Atlasfile_ServiceDesc is the grpc.ServiceDesc for Atlasfile service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Atlasfile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sdk.Atlasfile",
	HandlerType: (*AtlasfileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Eval",
			Handler:    _Atlasfile_Eval_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Atlasfile_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sdk.proto",
}
