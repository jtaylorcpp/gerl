// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc.proto

package core

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GerlMsg_Type int32

const (
	GerlMsg_CALL GerlMsg_Type = 0
	GerlMsg_CAST GerlMsg_Type = 1
	GerlMsg_PROC GerlMsg_Type = 2
)

var GerlMsg_Type_name = map[int32]string{
	0: "CALL",
	1: "CAST",
	2: "PROC",
}

var GerlMsg_Type_value = map[string]int32{
	"CALL": 0,
	"CAST": 1,
	"PROC": 2,
}

func (x GerlMsg_Type) String() string {
	return proto.EnumName(GerlMsg_Type_name, int32(x))
}

func (GerlMsg_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{1, 0}
}

type Message_Type int32

const (
	Message_SIMPLE Message_Type = 0
)

var Message_Type_name = map[int32]string{
	0: "SIMPLE",
}

var Message_Type_value = map[string]int32{
	"SIMPLE": 0,
}

func (x Message_Type) String() string {
	return proto.EnumName(Message_Type_name, int32(x))
}

func (Message_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{2, 0}
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type GerlMsg struct {
	Type                 GerlMsg_Type `protobuf:"varint,1,opt,name=type,proto3,enum=core.GerlMsg_Type" json:"type,omitempty"`
	Fromaddr             string       `protobuf:"bytes,2,opt,name=fromaddr,proto3" json:"fromaddr,omitempty"`
	Msg                  *Message     `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *GerlMsg) Reset()         { *m = GerlMsg{} }
func (m *GerlMsg) String() string { return proto.CompactTextString(m) }
func (*GerlMsg) ProtoMessage()    {}
func (*GerlMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{1}
}

func (m *GerlMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GerlMsg.Unmarshal(m, b)
}
func (m *GerlMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GerlMsg.Marshal(b, m, deterministic)
}
func (m *GerlMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GerlMsg.Merge(m, src)
}
func (m *GerlMsg) XXX_Size() int {
	return xxx_messageInfo_GerlMsg.Size(m)
}
func (m *GerlMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_GerlMsg.DiscardUnknown(m)
}

var xxx_messageInfo_GerlMsg proto.InternalMessageInfo

func (m *GerlMsg) GetType() GerlMsg_Type {
	if m != nil {
		return m.Type
	}
	return GerlMsg_CALL
}

func (m *GerlMsg) GetFromaddr() string {
	if m != nil {
		return m.Fromaddr
	}
	return ""
}

func (m *GerlMsg) GetMsg() *Message {
	if m != nil {
		return m.Msg
	}
	return nil
}

type Message struct {
	Type                 Message_Type `protobuf:"varint,1,opt,name=type,proto3,enum=core.Message_Type" json:"type,omitempty"`
	Description          string       `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{2}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetType() Message_Type {
	if m != nil {
		return m.Type
	}
	return Message_SIMPLE
}

func (m *Message) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterEnum("core.GerlMsg_Type", GerlMsg_Type_name, GerlMsg_Type_value)
	proto.RegisterEnum("core.Message_Type", Message_Type_name, Message_Type_value)
	proto.RegisterType((*Empty)(nil), "core.Empty")
	proto.RegisterType((*GerlMsg)(nil), "core.GerlMsg")
	proto.RegisterType((*Message)(nil), "core.Message")
}

func init() { proto.RegisterFile("grpc.proto", fileDescriptor_bedfbfc9b54e5600) }

var fileDescriptor_bedfbfc9b54e5600 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x50, 0xc1, 0x6a, 0x83, 0x40,
	0x14, 0x74, 0x13, 0x1b, 0xd3, 0x67, 0x53, 0xe4, 0x9d, 0x24, 0x97, 0xca, 0x12, 0x8a, 0x27, 0x0f,
	0xf6, 0x0b, 0x8a, 0x84, 0x52, 0x50, 0x1a, 0x4c, 0x8e, 0xbd, 0x58, 0xdd, 0x4a, 0x40, 0xbb, 0xcb,
	0xee, 0x5e, 0xfc, 0x8f, 0x7e, 0x70, 0x59, 0x77, 0x29, 0x49, 0xe9, 0xed, 0xcd, 0xce, 0xec, 0xbc,
	0x79, 0x03, 0xd0, 0x4b, 0xd1, 0x66, 0x42, 0x72, 0xcd, 0xd1, 0x6f, 0xb9, 0x64, 0x34, 0x80, 0x9b,
	0xfd, 0x28, 0xf4, 0x44, 0xbf, 0x09, 0x04, 0x2f, 0x4c, 0x0e, 0x95, 0xea, 0xf1, 0x11, 0x7c, 0x3d,
	0x09, 0x16, 0x93, 0x84, 0xa4, 0xf7, 0x39, 0x66, 0x46, 0x99, 0x39, 0x32, 0x3b, 0x4d, 0x82, 0xd5,
	0x33, 0x8f, 0x5b, 0x58, 0x7f, 0x4a, 0x3e, 0x36, 0x5d, 0x27, 0xe3, 0x45, 0x42, 0xd2, 0xdb, 0xfa,
	0x17, 0xe3, 0x03, 0x2c, 0x47, 0xd5, 0xc7, 0xcb, 0x84, 0xa4, 0x61, 0xbe, 0xb1, 0x16, 0x15, 0x53,
	0xaa, 0xe9, 0x59, 0x6d, 0x18, 0xba, 0x03, 0xdf, 0x58, 0xe1, 0x1a, 0xfc, 0xe2, 0xb9, 0x2c, 0x23,
	0xcf, 0x4e, 0xc7, 0x53, 0x44, 0xcc, 0x74, 0xa8, 0xdf, 0x8a, 0x68, 0x41, 0x7b, 0x08, 0xdc, 0xaf,
	0xff, 0x53, 0x39, 0xf2, 0x32, 0x55, 0x02, 0x61, 0xc7, 0x54, 0x2b, 0xcf, 0x42, 0x9f, 0xf9, 0x97,
	0x0b, 0x76, 0xf9, 0x44, 0xd1, 0xad, 0x06, 0x58, 0x1d, 0x5f, 0xab, 0x43, 0xb9, 0x8f, 0xbc, 0xfc,
	0x1d, 0xee, 0xe6, 0x0b, 0xad, 0x9f, 0x34, 0xdb, 0x8a, 0x66, 0x18, 0x70, 0x73, 0x75, 0xfd, 0xf6,
	0x1a, 0x52, 0x0f, 0x77, 0x46, 0xa7, 0xf4, 0x5f, 0x5d, 0x68, 0xa1, 0xed, 0xd6, 0xfb, 0x58, 0xcd,
	0x9d, 0x3f, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0x94, 0x74, 0x84, 0xfa, 0x81, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GerlMessagerClient is the client API for GerlMessager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GerlMessagerClient interface {
	Call(ctx context.Context, in *GerlMsg, opts ...grpc.CallOption) (*GerlMsg, error)
	Cast(ctx context.Context, in *GerlMsg, opts ...grpc.CallOption) (*Empty, error)
}

type gerlMessagerClient struct {
	cc *grpc.ClientConn
}

func NewGerlMessagerClient(cc *grpc.ClientConn) GerlMessagerClient {
	return &gerlMessagerClient{cc}
}

func (c *gerlMessagerClient) Call(ctx context.Context, in *GerlMsg, opts ...grpc.CallOption) (*GerlMsg, error) {
	out := new(GerlMsg)
	err := c.cc.Invoke(ctx, "/core.GerlMessager/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gerlMessagerClient) Cast(ctx context.Context, in *GerlMsg, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/core.GerlMessager/Cast", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GerlMessagerServer is the server API for GerlMessager service.
type GerlMessagerServer interface {
	Call(context.Context, *GerlMsg) (*GerlMsg, error)
	Cast(context.Context, *GerlMsg) (*Empty, error)
}

func RegisterGerlMessagerServer(s *grpc.Server, srv GerlMessagerServer) {
	s.RegisterService(&_GerlMessager_serviceDesc, srv)
}

func _GerlMessager_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GerlMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GerlMessagerServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.GerlMessager/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GerlMessagerServer).Call(ctx, req.(*GerlMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _GerlMessager_Cast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GerlMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GerlMessagerServer).Cast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.GerlMessager/Cast",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GerlMessagerServer).Cast(ctx, req.(*GerlMsg))
	}
	return interceptor(ctx, in, info, handler)
}

var _GerlMessager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "core.GerlMessager",
	HandlerType: (*GerlMessagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _GerlMessager_Call_Handler,
		},
		{
			MethodName: "Cast",
			Handler:    _GerlMessager_Cast_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc.proto",
}
