// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: proto/biz.proto

package __

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
	BizService_GetUsers_FullMethodName        = "/BizService/GetUsers"
	BizService_GetUsersWithSQL_FullMethodName = "/BizService/GetUsersWithSQL"
)

// BizServiceClient is the client API for BizService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BizServiceClient interface {
	GetUsers(ctx context.Context, in *GetUsersRequest, opts ...grpc.CallOption) (*GetUsersResponse, error)
	GetUsersWithSQL(ctx context.Context, in *GetUsersWithSQLRequest, opts ...grpc.CallOption) (*GetUsersResponse, error)
}

type bizServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBizServiceClient(cc grpc.ClientConnInterface) BizServiceClient {
	return &bizServiceClient{cc}
}

func (c *bizServiceClient) GetUsers(ctx context.Context, in *GetUsersRequest, opts ...grpc.CallOption) (*GetUsersResponse, error) {
	out := new(GetUsersResponse)
	err := c.cc.Invoke(ctx, BizService_GetUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bizServiceClient) GetUsersWithSQL(ctx context.Context, in *GetUsersWithSQLRequest, opts ...grpc.CallOption) (*GetUsersResponse, error) {
	out := new(GetUsersResponse)
	err := c.cc.Invoke(ctx, BizService_GetUsersWithSQL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BizServiceServer is the server API for BizService service.
// All implementations must embed UnimplementedBizServiceServer
// for forward compatibility
type BizServiceServer interface {
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)
	GetUsersWithSQL(context.Context, *GetUsersWithSQLRequest) (*GetUsersResponse, error)
	mustEmbedUnimplementedBizServiceServer()
}

// UnimplementedBizServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBizServiceServer struct {
}

func (UnimplementedBizServiceServer) GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}
func (UnimplementedBizServiceServer) GetUsersWithSQL(context.Context, *GetUsersWithSQLRequest) (*GetUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsersWithSQL not implemented")
}
func (UnimplementedBizServiceServer) mustEmbedUnimplementedBizServiceServer() {}

// UnsafeBizServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BizServiceServer will
// result in compilation errors.
type UnsafeBizServiceServer interface {
	mustEmbedUnimplementedBizServiceServer()
}

func RegisterBizServiceServer(s grpc.ServiceRegistrar, srv BizServiceServer) {
	s.RegisterService(&BizService_ServiceDesc, srv)
}

func _BizService_GetUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BizServiceServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BizService_GetUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BizServiceServer).GetUsers(ctx, req.(*GetUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BizService_GetUsersWithSQL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersWithSQLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BizServiceServer).GetUsersWithSQL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BizService_GetUsersWithSQL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BizServiceServer).GetUsersWithSQL(ctx, req.(*GetUsersWithSQLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BizService_ServiceDesc is the grpc.ServiceDesc for BizService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BizService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "BizService",
	HandlerType: (*BizServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUsers",
			Handler:    _BizService_GetUsers_Handler,
		},
		{
			MethodName: "GetUsersWithSQL",
			Handler:    _BizService_GetUsersWithSQL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/biz.proto",
}
