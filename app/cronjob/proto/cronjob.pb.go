// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cronjob.proto

package cronjob

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import protocol "ashe/protocol"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CronjobClient is the client API for Cronjob service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CronjobClient interface {
	Add(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error)
	//    rpc edit (protocol.Request) returns (protocol.Response) {
	//    };
	//    rpc search (protocol.Request) returns (protocol.Response) {
	//    };
	Del(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error)
	List(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error)
	Sync(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error)
}

type cronjobClient struct {
	cc *grpc.ClientConn
}

func NewCronjobClient(cc *grpc.ClientConn) CronjobClient {
	return &cronjobClient{cc}
}

func (c *cronjobClient) Add(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
	out := new(protocol.Response)
	err := c.cc.Invoke(ctx, "/cronjob/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cronjobClient) Del(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
	out := new(protocol.Response)
	err := c.cc.Invoke(ctx, "/cronjob/Del", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cronjobClient) List(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
	out := new(protocol.Response)
	err := c.cc.Invoke(ctx, "/cronjob/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cronjobClient) Sync(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
	out := new(protocol.Response)
	err := c.cc.Invoke(ctx, "/cronjob/Sync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CronjobServer is the server API for Cronjob service.
type CronjobServer interface {
	Add(context.Context, *protocol.Request) (*protocol.Response, error)
	//    rpc edit (protocol.Request) returns (protocol.Response) {
	//    };
	//    rpc search (protocol.Request) returns (protocol.Response) {
	//    };
	Del(context.Context, *protocol.Request) (*protocol.Response, error)
	List(context.Context, *protocol.Request) (*protocol.Response, error)
	Sync(context.Context, *protocol.Request) (*protocol.Response, error)
}

func RegisterCronjobServer(s *grpc.Server, srv CronjobServer) {
	s.RegisterService(&_Cronjob_serviceDesc, srv)
}

func _Cronjob_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protocol.Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronjobServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cronjob/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronjobServer).Add(ctx, req.(*protocol.Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cronjob_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protocol.Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronjobServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cronjob/Del",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronjobServer).Del(ctx, req.(*protocol.Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cronjob_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protocol.Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronjobServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cronjob/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronjobServer).List(ctx, req.(*protocol.Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cronjob_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protocol.Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronjobServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cronjob/Sync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronjobServer).Sync(ctx, req.(*protocol.Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Cronjob_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cronjob",
	HandlerType: (*CronjobServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Cronjob_Add_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _Cronjob_Del_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Cronjob_List_Handler,
		},
		{
			MethodName: "Sync",
			Handler:    _Cronjob_Sync_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cronjob.proto",
}

func init() { proto.RegisterFile("cronjob.proto", fileDescriptor_cronjob_cf9f1248f5a26b9a) }

var fileDescriptor_cronjob_cf9f1248f5a26b9a = []byte{
	// 121 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0x2e, 0xca, 0xcf,
	0xcb, 0xca, 0x4f, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0x92, 0x49, 0x2c, 0xce, 0x48, 0xd5,
	0x07, 0xb3, 0x93, 0xf3, 0x73, 0xe0, 0x0c, 0x88, 0xac, 0xd1, 0x69, 0x46, 0x2e, 0x76, 0xa8, 0x7a,
	0x21, 0x3d, 0x2e, 0x66, 0xc7, 0x94, 0x14, 0x21, 0x41, 0x3d, 0xb8, 0x9a, 0xa0, 0xd4, 0xc2, 0xd2,
	0xd4, 0xe2, 0x12, 0x29, 0x21, 0x64, 0xa1, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x25, 0x06, 0x90,
	0x7a, 0x97, 0xd4, 0x1c, 0xe2, 0xd5, 0xeb, 0x73, 0xb1, 0xf8, 0x64, 0x16, 0x97, 0x90, 0xa4, 0x21,
	0xb8, 0x32, 0x2f, 0x99, 0x68, 0x0d, 0x49, 0x6c, 0x60, 0x41, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xd2, 0x21, 0xff, 0xa8, 0x03, 0x01, 0x00, 0x00,
}
