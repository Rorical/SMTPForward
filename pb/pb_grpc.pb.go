// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: pb.proto

package pb

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
	SMTPForward_SendSMTP_FullMethodName = "/pb.SMTPForward/SendSMTP"
)

// SMTPForwardClient is the client API for SMTPForward service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SMTPForwardClient interface {
	SendSMTP(ctx context.Context, in *SMTPData, opts ...grpc.CallOption) (*SMTPResult, error)
}

type sMTPForwardClient struct {
	cc grpc.ClientConnInterface
}

func NewSMTPForwardClient(cc grpc.ClientConnInterface) SMTPForwardClient {
	return &sMTPForwardClient{cc}
}

func (c *sMTPForwardClient) SendSMTP(ctx context.Context, in *SMTPData, opts ...grpc.CallOption) (*SMTPResult, error) {
	out := new(SMTPResult)
	err := c.cc.Invoke(ctx, SMTPForward_SendSMTP_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SMTPForwardServer is the server API for SMTPForward service.
// All implementations must embed UnimplementedSMTPForwardServer
// for forward compatibility
type SMTPForwardServer interface {
	SendSMTP(context.Context, *SMTPData) (*SMTPResult, error)
	mustEmbedUnimplementedSMTPForwardServer()
}

// UnimplementedSMTPForwardServer must be embedded to have forward compatible implementations.
type UnimplementedSMTPForwardServer struct {
}

func (UnimplementedSMTPForwardServer) SendSMTP(context.Context, *SMTPData) (*SMTPResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSMTP not implemented")
}
func (UnimplementedSMTPForwardServer) mustEmbedUnimplementedSMTPForwardServer() {}

// UnsafeSMTPForwardServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SMTPForwardServer will
// result in compilation errors.
type UnsafeSMTPForwardServer interface {
	mustEmbedUnimplementedSMTPForwardServer()
}

func RegisterSMTPForwardServer(s grpc.ServiceRegistrar, srv SMTPForwardServer) {
	s.RegisterService(&SMTPForward_ServiceDesc, srv)
}

func _SMTPForward_SendSMTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SMTPData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SMTPForwardServer).SendSMTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SMTPForward_SendSMTP_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SMTPForwardServer).SendSMTP(ctx, req.(*SMTPData))
	}
	return interceptor(ctx, in, info, handler)
}

// SMTPForward_ServiceDesc is the grpc.ServiceDesc for SMTPForward service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SMTPForward_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.SMTPForward",
	HandlerType: (*SMTPForwardServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendSMTP",
			Handler:    _SMTPForward_SendSMTP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb.proto",
}