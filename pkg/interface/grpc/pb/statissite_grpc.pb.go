// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// StaticSiteServiceClient is the client API for StaticSiteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StaticSiteServiceClient interface {
	Reload(ctx context.Context, in *ReloadRequest, opts ...grpc.CallOption) (*ReloadResponse, error)
}

type staticSiteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStaticSiteServiceClient(cc grpc.ClientConnInterface) StaticSiteServiceClient {
	return &staticSiteServiceClient{cc}
}

func (c *staticSiteServiceClient) Reload(ctx context.Context, in *ReloadRequest, opts ...grpc.CallOption) (*ReloadResponse, error) {
	out := new(ReloadResponse)
	err := c.cc.Invoke(ctx, "/neoshowcase.protobuf.StaticSiteService/Reload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StaticSiteServiceServer is the server API for StaticSiteService service.
// All implementations must embed UnimplementedStaticSiteServiceServer
// for forward compatibility
type StaticSiteServiceServer interface {
	Reload(context.Context, *ReloadRequest) (*ReloadResponse, error)
	mustEmbedUnimplementedStaticSiteServiceServer()
}

// UnimplementedStaticSiteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStaticSiteServiceServer struct {
}

func (UnimplementedStaticSiteServiceServer) Reload(context.Context, *ReloadRequest) (*ReloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reload not implemented")
}
func (UnimplementedStaticSiteServiceServer) mustEmbedUnimplementedStaticSiteServiceServer() {}

// UnsafeStaticSiteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StaticSiteServiceServer will
// result in compilation errors.
type UnsafeStaticSiteServiceServer interface {
	mustEmbedUnimplementedStaticSiteServiceServer()
}

func RegisterStaticSiteServiceServer(s grpc.ServiceRegistrar, srv StaticSiteServiceServer) {
	s.RegisterService(&_StaticSiteService_serviceDesc, srv)
}

func _StaticSiteService_Reload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaticSiteServiceServer).Reload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neoshowcase.protobuf.StaticSiteService/Reload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaticSiteServiceServer).Reload(ctx, req.(*ReloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _StaticSiteService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "neoshowcase.protobuf.StaticSiteService",
	HandlerType: (*StaticSiteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Reload",
			Handler:    _StaticSiteService_Reload_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "neoshowcase/protobuf/statissite.proto",
}
