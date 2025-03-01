//protoc -I protobuf --go_out=paths=source_relative:internal/api/gen --go-grpc_out=paths=source_relative:internal/api/gen protobuf/api.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: api.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UrlShortener_RegUrl_FullMethodName = "/api.UrlShortener/RegUrl"
	UrlShortener_GetUrl_FullMethodName = "/api.UrlShortener/GetUrl"
)

// UrlShortenerClient is the client API for UrlShortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UrlShortenerClient interface {
	RegUrl(ctx context.Context, in *RegUrlReq, opts ...grpc.CallOption) (*RegUrlResp, error)
	GetUrl(ctx context.Context, in *GetUrlReq, opts ...grpc.CallOption) (*GetUrlResp, error)
}

type urlShortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewUrlShortenerClient(cc grpc.ClientConnInterface) UrlShortenerClient {
	return &urlShortenerClient{cc}
}

func (c *urlShortenerClient) RegUrl(ctx context.Context, in *RegUrlReq, opts ...grpc.CallOption) (*RegUrlResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegUrlResp)
	err := c.cc.Invoke(ctx, UrlShortener_RegUrl_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *urlShortenerClient) GetUrl(ctx context.Context, in *GetUrlReq, opts ...grpc.CallOption) (*GetUrlResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUrlResp)
	err := c.cc.Invoke(ctx, UrlShortener_GetUrl_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UrlShortenerServer is the server API for UrlShortener service.
// All implementations must embed UnimplementedUrlShortenerServer
// for forward compatibility.
type UrlShortenerServer interface {
	RegUrl(context.Context, *RegUrlReq) (*RegUrlResp, error)
	GetUrl(context.Context, *GetUrlReq) (*GetUrlResp, error)
	mustEmbedUnimplementedUrlShortenerServer()
}

// UnimplementedUrlShortenerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUrlShortenerServer struct{}

func (UnimplementedUrlShortenerServer) RegUrl(context.Context, *RegUrlReq) (*RegUrlResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegUrl not implemented")
}
func (UnimplementedUrlShortenerServer) GetUrl(context.Context, *GetUrlReq) (*GetUrlResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUrl not implemented")
}
func (UnimplementedUrlShortenerServer) mustEmbedUnimplementedUrlShortenerServer() {}
func (UnimplementedUrlShortenerServer) testEmbeddedByValue()                      {}

// UnsafeUrlShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UrlShortenerServer will
// result in compilation errors.
type UnsafeUrlShortenerServer interface {
	mustEmbedUnimplementedUrlShortenerServer()
}

func RegisterUrlShortenerServer(s grpc.ServiceRegistrar, srv UrlShortenerServer) {
	// If the following call pancis, it indicates UnimplementedUrlShortenerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UrlShortener_ServiceDesc, srv)
}

func _UrlShortener_RegUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegUrlReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlShortenerServer).RegUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UrlShortener_RegUrl_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlShortenerServer).RegUrl(ctx, req.(*RegUrlReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UrlShortener_GetUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUrlReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlShortenerServer).GetUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UrlShortener_GetUrl_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlShortenerServer).GetUrl(ctx, req.(*GetUrlReq))
	}
	return interceptor(ctx, in, info, handler)
}

// UrlShortener_ServiceDesc is the grpc.ServiceDesc for UrlShortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UrlShortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.UrlShortener",
	HandlerType: (*UrlShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegUrl",
			Handler:    _UrlShortener_RegUrl_Handler,
		},
		{
			MethodName: "GetUrl",
			Handler:    _UrlShortener_GetUrl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
