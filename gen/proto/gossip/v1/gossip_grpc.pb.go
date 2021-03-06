// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: proto/gossip/v1/gossip.proto

package gossipv1

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

// GossipServiceClient is the client API for GossipService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GossipServiceClient interface {
	// rpc Increment(IncrementRequest) returns (IncrementResponse) {}
	Merge(ctx context.Context, in *MergeRequest, opts ...grpc.CallOption) (*MergeResponse, error)
	Ratelimit(ctx context.Context, in *RatelimitRequest, opts ...grpc.CallOption) (*RatelimitResponse, error)
}

type gossipServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGossipServiceClient(cc grpc.ClientConnInterface) GossipServiceClient {
	return &gossipServiceClient{cc}
}

func (c *gossipServiceClient) Merge(ctx context.Context, in *MergeRequest, opts ...grpc.CallOption) (*MergeResponse, error) {
	out := new(MergeResponse)
	err := c.cc.Invoke(ctx, "/proto.gossip.v1.GossipService/Merge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gossipServiceClient) Ratelimit(ctx context.Context, in *RatelimitRequest, opts ...grpc.CallOption) (*RatelimitResponse, error) {
	out := new(RatelimitResponse)
	err := c.cc.Invoke(ctx, "/proto.gossip.v1.GossipService/Ratelimit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GossipServiceServer is the server API for GossipService service.
// All implementations must embed UnimplementedGossipServiceServer
// for forward compatibility
type GossipServiceServer interface {
	// rpc Increment(IncrementRequest) returns (IncrementResponse) {}
	Merge(context.Context, *MergeRequest) (*MergeResponse, error)
	Ratelimit(context.Context, *RatelimitRequest) (*RatelimitResponse, error)
	mustEmbedUnimplementedGossipServiceServer()
}

// UnimplementedGossipServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGossipServiceServer struct {
}

func (UnimplementedGossipServiceServer) Merge(context.Context, *MergeRequest) (*MergeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Merge not implemented")
}
func (UnimplementedGossipServiceServer) Ratelimit(context.Context, *RatelimitRequest) (*RatelimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ratelimit not implemented")
}
func (UnimplementedGossipServiceServer) mustEmbedUnimplementedGossipServiceServer() {}

// UnsafeGossipServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GossipServiceServer will
// result in compilation errors.
type UnsafeGossipServiceServer interface {
	mustEmbedUnimplementedGossipServiceServer()
}

func RegisterGossipServiceServer(s grpc.ServiceRegistrar, srv GossipServiceServer) {
	s.RegisterService(&GossipService_ServiceDesc, srv)
}

func _GossipService_Merge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GossipServiceServer).Merge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gossip.v1.GossipService/Merge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GossipServiceServer).Merge(ctx, req.(*MergeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GossipService_Ratelimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RatelimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GossipServiceServer).Ratelimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gossip.v1.GossipService/Ratelimit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GossipServiceServer).Ratelimit(ctx, req.(*RatelimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GossipService_ServiceDesc is the grpc.ServiceDesc for GossipService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GossipService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.gossip.v1.GossipService",
	HandlerType: (*GossipServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Merge",
			Handler:    _GossipService_Merge_Handler,
		},
		{
			MethodName: "Ratelimit",
			Handler:    _GossipService_Ratelimit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/gossip/v1/gossip.proto",
}
