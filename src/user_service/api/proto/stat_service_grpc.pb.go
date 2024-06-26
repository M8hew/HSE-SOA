// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: stat_service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StatServiceClient is the client API for StatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatServiceClient interface {
	GetViewsLikes(ctx context.Context, in *GetViewsLikesRequest, opts ...grpc.CallOption) (*GetViewsLikesResponse, error)
	GetTopPosts(ctx context.Context, in *GetTopPostsRequest, opts ...grpc.CallOption) (*GetTopPostsResponse, error)
	GetTopUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetTopUsersResponse, error)
}

type statServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStatServiceClient(cc grpc.ClientConnInterface) StatServiceClient {
	return &statServiceClient{cc}
}

func (c *statServiceClient) GetViewsLikes(ctx context.Context, in *GetViewsLikesRequest, opts ...grpc.CallOption) (*GetViewsLikesResponse, error) {
	out := new(GetViewsLikesResponse)
	err := c.cc.Invoke(ctx, "/stats.StatService/GetViewsLikes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statServiceClient) GetTopPosts(ctx context.Context, in *GetTopPostsRequest, opts ...grpc.CallOption) (*GetTopPostsResponse, error) {
	out := new(GetTopPostsResponse)
	err := c.cc.Invoke(ctx, "/stats.StatService/GetTopPosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statServiceClient) GetTopUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetTopUsersResponse, error) {
	out := new(GetTopUsersResponse)
	err := c.cc.Invoke(ctx, "/stats.StatService/GetTopUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatServiceServer is the server API for StatService service.
// All implementations must embed UnimplementedStatServiceServer
// for forward compatibility
type StatServiceServer interface {
	GetViewsLikes(context.Context, *GetViewsLikesRequest) (*GetViewsLikesResponse, error)
	GetTopPosts(context.Context, *GetTopPostsRequest) (*GetTopPostsResponse, error)
	GetTopUsers(context.Context, *emptypb.Empty) (*GetTopUsersResponse, error)
	mustEmbedUnimplementedStatServiceServer()
}

// UnimplementedStatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStatServiceServer struct {
}

func (UnimplementedStatServiceServer) GetViewsLikes(context.Context, *GetViewsLikesRequest) (*GetViewsLikesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetViewsLikes not implemented")
}
func (UnimplementedStatServiceServer) GetTopPosts(context.Context, *GetTopPostsRequest) (*GetTopPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopPosts not implemented")
}
func (UnimplementedStatServiceServer) GetTopUsers(context.Context, *emptypb.Empty) (*GetTopUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopUsers not implemented")
}
func (UnimplementedStatServiceServer) mustEmbedUnimplementedStatServiceServer() {}

// UnsafeStatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatServiceServer will
// result in compilation errors.
type UnsafeStatServiceServer interface {
	mustEmbedUnimplementedStatServiceServer()
}

func RegisterStatServiceServer(s grpc.ServiceRegistrar, srv StatServiceServer) {
	s.RegisterService(&StatService_ServiceDesc, srv)
}

func _StatService_GetViewsLikes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetViewsLikesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServiceServer).GetViewsLikes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.StatService/GetViewsLikes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServiceServer).GetViewsLikes(ctx, req.(*GetViewsLikesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatService_GetTopPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTopPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServiceServer).GetTopPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.StatService/GetTopPosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServiceServer).GetTopPosts(ctx, req.(*GetTopPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatService_GetTopUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServiceServer).GetTopUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.StatService/GetTopUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServiceServer).GetTopUsers(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// StatService_ServiceDesc is the grpc.ServiceDesc for StatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stats.StatService",
	HandlerType: (*StatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetViewsLikes",
			Handler:    _StatService_GetViewsLikes_Handler,
		},
		{
			MethodName: "GetTopPosts",
			Handler:    _StatService_GetTopPosts_Handler,
		},
		{
			MethodName: "GetTopUsers",
			Handler:    _StatService_GetTopUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stat_service.proto",
}
