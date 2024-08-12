// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: tei.proto

package generated

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
	Rerank_Rerank_FullMethodName = "/tei.v1.Rerank/Rerank"
)

// RerankClient is the client API for Rerank service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RerankClient interface {
	Rerank(ctx context.Context, in *RerankRequest, opts ...grpc.CallOption) (*RerankResponse, error)
}

type rerankClient struct {
	cc grpc.ClientConnInterface
}

func NewRerankClient(cc grpc.ClientConnInterface) RerankClient {
	return &rerankClient{cc}
}

func (c *rerankClient) Rerank(ctx context.Context, in *RerankRequest, opts ...grpc.CallOption) (*RerankResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RerankResponse)
	err := c.cc.Invoke(ctx, Rerank_Rerank_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RerankServer is the server API for Rerank service.
// All implementations must embed UnimplementedRerankServer
// for forward compatibility.
type RerankServer interface {
	Rerank(context.Context, *RerankRequest) (*RerankResponse, error)
	mustEmbedUnimplementedRerankServer()
}

// UnimplementedRerankServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRerankServer struct{}

func (UnimplementedRerankServer) Rerank(context.Context, *RerankRequest) (*RerankResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rerank not implemented")
}
func (UnimplementedRerankServer) mustEmbedUnimplementedRerankServer() {}
func (UnimplementedRerankServer) testEmbeddedByValue()                {}

// UnsafeRerankServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RerankServer will
// result in compilation errors.
type UnsafeRerankServer interface {
	mustEmbedUnimplementedRerankServer()
}

func RegisterRerankServer(s grpc.ServiceRegistrar, srv RerankServer) {
	// If the following call pancis, it indicates UnimplementedRerankServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Rerank_ServiceDesc, srv)
}

func _Rerank_Rerank_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RerankRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RerankServer).Rerank(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rerank_Rerank_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RerankServer).Rerank(ctx, req.(*RerankRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Rerank_ServiceDesc is the grpc.ServiceDesc for Rerank service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rerank_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tei.v1.Rerank",
	HandlerType: (*RerankServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Rerank",
			Handler:    _Rerank_Rerank_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tei.proto",
}