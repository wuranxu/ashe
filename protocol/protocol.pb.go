// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol.proto

package protocol

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Request struct {
	RequestJson          string   `protobuf:"bytes,1,opt,name=requestJson,proto3" json:"requestJson,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetRequestJson() string {
	if m != nil {
		return m.RequestJson
	}
	return ""
}

type Response struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	ResultJson           string   `protobuf:"bytes,3,opt,name=resultJson,proto3" json:"resultJson,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
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

func (m *Response) GetResultJson() string {
	if m != nil {
		return m.ResultJson
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "Request")
	proto.RegisterType((*Response)(nil), "Response")
}

func init() { proto.RegisterFile("protocol.proto", fileDescriptor_2bc2336598a3f7e0) }

var fileDescriptor_2bc2336598a3f7e0 = []byte{
	// 127 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x03, 0x33, 0x94, 0xb4, 0xb9, 0xd8, 0x83, 0x52, 0x0b, 0x4b, 0x53,
	0x8b, 0x4b, 0x84, 0x14, 0xb8, 0xb8, 0x8b, 0x20, 0x4c, 0xaf, 0xe2, 0xfc, 0x3c, 0x09, 0x46, 0x05,
	0x46, 0x0d, 0xce, 0x20, 0x64, 0x21, 0xa5, 0x00, 0x2e, 0x8e, 0xa0, 0xd4, 0xe2, 0x82, 0xfc, 0xbc,
	0xe2, 0x54, 0x21, 0x21, 0x2e, 0x96, 0xe4, 0xfc, 0x94, 0x54, 0xb0, 0x32, 0xd6, 0x20, 0x30, 0x5b,
	0x48, 0x80, 0x8b, 0x39, 0xb7, 0x38, 0x5d, 0x82, 0x09, 0xac, 0x13, 0xc4, 0x14, 0x92, 0xe3, 0xe2,
	0x2a, 0x4a, 0x2d, 0x2e, 0xcd, 0x81, 0x18, 0xc9, 0x0c, 0x96, 0x40, 0x12, 0x49, 0x62, 0x03, 0xbb,
	0xc2, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xd3, 0xba, 0xf4, 0x5f, 0x97, 0x00, 0x00, 0x00,
}
