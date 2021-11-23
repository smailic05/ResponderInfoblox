// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// MyResponderClient is the client API for MyResponder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MyResponderClient interface {
	GetVersion(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*VersionResponse, error)
}

type myResponderClient struct {
	cc grpc.ClientConnInterface
}

func NewMyResponderClient(cc grpc.ClientConnInterface) MyResponderClient {
	return &myResponderClient{cc}
}

func (c *myResponderClient) GetVersion(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/myresponder.MyResponder/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MyResponderServer is the server API for MyResponder service.
// All implementations should embed UnimplementedMyResponderServer
// for forward compatibility
type MyResponderServer interface {
	GetVersion(context.Context, *emptypb.Empty) (*VersionResponse, error)
}

// UnimplementedMyResponderServer should be embedded to have forward compatible implementations.
type UnimplementedMyResponderServer struct {
}

func (UnimplementedMyResponderServer) GetVersion(context.Context, *emptypb.Empty) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}

// UnsafeMyResponderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MyResponderServer will
// result in compilation errors.
type UnsafeMyResponderServer interface {
	mustEmbedUnimplementedMyResponderServer()
}

func RegisterMyResponderServer(s grpc.ServiceRegistrar, srv MyResponderServer) {
	s.RegisterService(&MyResponder_ServiceDesc, srv)
}

func _MyResponder_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MyResponderServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/myresponder.MyResponder/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MyResponderServer).GetVersion(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MyResponder_ServiceDesc is the grpc.ServiceDesc for MyResponder service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MyResponder_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "myresponder.MyResponder",
	HandlerType: (*MyResponderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _MyResponder_GetVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/my-responder/pkg/pb/service.proto",
}
