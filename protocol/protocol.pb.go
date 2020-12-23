// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol.proto

package protocol // import "ashe/protocol"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

type Request struct {
	RequestJson          []byte   `protobuf:"bytes,1,opt,name=requestJson,proto3" json:"requestJson,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_47966471535ff2c1, []int{0}
}
func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (dst *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(dst, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetRequestJson() []byte {
	if m != nil {
		return m.RequestJson
	}
	return nil
}

type Response struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	ResultJson           []byte   `protobuf:"bytes,3,opt,name=resultJson,proto3" json:"resultJson,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_47966471535ff2c1, []int{1}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Response) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *Response) GetResultJson() []byte {
	if m != nil {
		return m.ResultJson
	}
	return nil
}

type Args struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Service              string   `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	Method               string   `protobuf:"bytes,3,opt,name=method,proto3" json:"method,omitempty"`
	Args                 []byte   `protobuf:"bytes,4,opt,name=args,proto3" json:"args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Args) Reset()         { *m = Args{} }
func (m *Args) String() string { return proto.CompactTextString(m) }
func (*Args) ProtoMessage()    {}
func (*Args) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_47966471535ff2c1, []int{2}
}
func (m *Args) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Args.Unmarshal(m, b)
}
func (m *Args) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Args.Marshal(b, m, deterministic)
}
func (dst *Args) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Args.Merge(dst, src)
}
func (m *Args) XXX_Size() int {
	return xxx_messageInfo_Args.Size(m)
}
func (m *Args) XXX_DiscardUnknown() {
	xxx_messageInfo_Args.DiscardUnknown(m)
}

var xxx_messageInfo_Args proto.InternalMessageInfo

func (m *Args) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Args) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *Args) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Args) GetArgs() []byte {
	if m != nil {
		return m.Args
	}
	return nil
}

func init() {
	proto.RegisterType((*Request)(nil), "Request")
	proto.RegisterType((*Response)(nil), "Response")
	proto.RegisterType((*Args)(nil), "Args")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RpcServiceClient is the client API for RpcService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RpcServiceClient interface {
	Invoke(ctx context.Context, in *Args, opts ...grpc.CallOption) (*Response, error)
}

type rpcServiceClient struct {
	cc *grpc.ClientConn
}

func NewRpcServiceClient(cc *grpc.ClientConn) RpcServiceClient {
	return &rpcServiceClient{cc}
}

func (c *rpcServiceClient) Invoke(ctx context.Context, in *Args, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/RpcService/Invoke", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RpcServiceServer is the server API for RpcService service.
type RpcServiceServer interface {
	Invoke(context.Context, *Args) (*Response, error)
}

func RegisterRpcServiceServer(s *grpc.Server, srv RpcServiceServer) {
	s.RegisterService(&_RpcService_serviceDesc, srv)
}

func _RpcService_Invoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServiceServer).Invoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcService/Invoke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServiceServer).Invoke(ctx, req.(*Args))
	}
	return interceptor(ctx, in, info, handler)
}

var _RpcService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "RpcService",
	HandlerType: (*RpcServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Invoke",
			Handler:    _RpcService_Invoke_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protocol.proto",
}

func init() { proto.RegisterFile("protocol.proto", fileDescriptor_protocol_47966471535ff2c1) }

var fileDescriptor_protocol_47966471535ff2c1 = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x4d, 0x4f, 0xc3, 0x30,
	0x0c, 0x86, 0x29, 0xeb, 0x3a, 0x6a, 0x3e, 0xe5, 0x03, 0xaa, 0x10, 0x42, 0x55, 0x4f, 0x13, 0x48,
	0x45, 0x82, 0x5f, 0x00, 0x37, 0x38, 0x21, 0x73, 0xe3, 0x36, 0x3a, 0xd3, 0x4d, 0x6c, 0x75, 0x89,
	0xbb, 0xfe, 0x7e, 0x14, 0x2f, 0x91, 0x76, 0x7b, 0x5e, 0x27, 0x79, 0xde, 0x24, 0x70, 0xd1, 0x3b,
	0x19, 0xa4, 0x91, 0x4d, 0x6d, 0x50, 0x3d, 0xc0, 0x8c, 0xf8, 0x6f, 0xc7, 0x3a, 0x60, 0x09, 0xa7,
	0x6e, 0x8f, 0xef, 0x2a, 0x5d, 0x91, 0x94, 0xc9, 0xfc, 0x8c, 0x0e, 0x47, 0xd5, 0x07, 0x9c, 0x10,
	0x6b, 0x2f, 0x9d, 0x32, 0x22, 0xa4, 0x8d, 0x2c, 0xd9, 0xb6, 0x4d, 0xc9, 0x18, 0xaf, 0x60, 0xb2,
	0xd5, 0xb6, 0x38, 0x2e, 0x93, 0x79, 0x4e, 0x1e, 0xf1, 0x0e, 0xc0, 0xb1, 0xee, 0x36, 0x7b, 0xe5,
	0xc4, 0x94, 0x07, 0x93, 0xea, 0x07, 0xd2, 0x17, 0xd7, 0x2a, 0x16, 0x30, 0x1b, 0xd9, 0xe9, 0x3a,
	0xf4, 0xe6, 0x14, 0xa3, 0x5f, 0x51, 0x76, 0xe3, 0xba, 0xe1, 0xe0, 0x8d, 0x11, 0xaf, 0x21, 0xdb,
	0xf2, 0xb0, 0x92, 0xa5, 0x79, 0x73, 0x0a, 0xc9, 0xdf, 0x6c, 0xe1, 0x5a, 0x2d, 0x52, 0x6b, 0x33,
	0x7e, 0xba, 0x07, 0xa0, 0xbe, 0xf9, 0x0c, 0x27, 0x6f, 0x21, 0x7b, 0xeb, 0x46, 0xf9, 0x65, 0x9c,
	0xd6, 0xbe, 0xfe, 0x26, 0xaf, 0xe3, 0xbb, 0xaa, 0xa3, 0xd7, 0xcb, 0xaf, 0xf3, 0x85, 0xae, 0xf8,
	0x31, 0xfe, 0xd4, 0x77, 0x66, 0xf4, 0xfc, 0x1f, 0x00, 0x00, 0xff, 0xff, 0xa5, 0x0a, 0x18, 0x62,
	0x3c, 0x01, 0x00, 0x00,
}
