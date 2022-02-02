// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.2
// source: v1/kgs.proto

package v1

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
	KGS_GetKeys_FullMethodName = "/v1.KGS/GetKeys"
)

// KGSClient is the client API for KGS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KGSClient interface {
	GetKeys(ctx context.Context, in *GetKeysRequest, opts ...grpc.CallOption) (*GetKeysReply, error)
}

type kGSClient struct {
	cc grpc.ClientConnInterface
}

func NewKGSClient(cc grpc.ClientConnInterface) KGSClient {
	return &kGSClient{cc}
}

func (c *kGSClient) GetKeys(ctx context.Context, in *GetKeysRequest, opts ...grpc.CallOption) (*GetKeysReply, error) {
	out := new(GetKeysReply)
	err := c.cc.Invoke(ctx, KGS_GetKeys_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KGSServer is the server API for KGS service.
// All implementations must embed UnimplementedKGSServer
// for forward compatibility
type KGSServer interface {
	GetKeys(context.Context, *GetKeysRequest) (*GetKeysReply, error)
	mustEmbedUnimplementedKGSServer()
}

// UnimplementedKGSServer must be embedded to have forward compatible implementations.
type UnimplementedKGSServer struct {
}

func (UnimplementedKGSServer) GetKeys(context.Context, *GetKeysRequest) (*GetKeysReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeys not implemented")
}
func (UnimplementedKGSServer) mustEmbedUnimplementedKGSServer() {}

// UnsafeKGSServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KGSServer will
// result in compilation errors.
type UnsafeKGSServer interface {
	mustEmbedUnimplementedKGSServer()
}

func RegisterKGSServer(s grpc.ServiceRegistrar, srv KGSServer) {
	s.RegisterService(&KGS_ServiceDesc, srv)
}

func _KGS_GetKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KGSServer).GetKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KGS_GetKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KGSServer).GetKeys(ctx, req.(*GetKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KGS_ServiceDesc is the grpc.ServiceDesc for KGS service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KGS_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.KGS",
	HandlerType: (*KGSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetKeys",
			Handler:    _KGS_GetKeys_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/kgs.proto",
}