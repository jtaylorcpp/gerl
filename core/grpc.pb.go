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
	Message_SIMPLE   Message_Type = 0
	Message_SYNC     Message_Type = 1
	Message_REGISTER Message_Type = 2
)

var Message_Type_name = map[int32]string{
	0: "SIMPLE",
	1: "SYNC",
	2: "REGISTER",
}

var Message_Type_value = map[string]int32{
	"SIMPLE":   0,
	"SYNC":     1,
	"REGISTER": 2,
}

func (x Message_Type) String() string {
	return proto.EnumName(Message_Type_name, int32(x))
}

func (Message_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{2, 0}
}

type Message_SubType int32

const (
	Message_GET  Message_SubType = 0
	Message_SET  Message_SubType = 1
	Message_PUT  Message_SubType = 2
	Message_JOIN Message_SubType = 3
)

var Message_SubType_name = map[int32]string{
	0: "GET",
	1: "SET",
	2: "PUT",
	3: "JOIN",
}

var Message_SubType_value = map[string]int32{
	"GET":  0,
	"SET":  1,
	"PUT":  2,
	"JOIN": 3,
}

func (x Message_SubType) String() string {
	return proto.EnumName(Message_SubType_name, int32(x))
}

func (Message_SubType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{2, 1}
}

type Health_Status int32

const (
	Health_ALIVE Health_Status = 0
)

var Health_Status_name = map[int32]string{
	0: "ALIVE",
}

var Health_Status_value = map[string]int32{
	"ALIVE": 0,
}

func (x Health_Status) String() string {
	return proto.EnumName(Health_Status_name, int32(x))
}

func (Health_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{3, 0}
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
	Type                 Message_Type    `protobuf:"varint,1,opt,name=type,proto3,enum=core.Message_Type" json:"type,omitempty"`
	Subtype              Message_SubType `protobuf:"varint,2,opt,name=subtype,proto3,enum=core.Message_SubType" json:"subtype,omitempty"`
	Description          string          `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Values               []string        `protobuf:"bytes,4,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
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

func (m *Message) GetSubtype() Message_SubType {
	if m != nil {
		return m.Subtype
	}
	return Message_GET
}

func (m *Message) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Message) GetValues() []string {
	if m != nil {
		return m.Values
	}
	return nil
}

type Health struct {
	Status               Health_Status `protobuf:"varint,1,opt,name=status,proto3,enum=core.Health_Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Health) Reset()         { *m = Health{} }
func (m *Health) String() string { return proto.CompactTextString(m) }
func (*Health) ProtoMessage()    {}
func (*Health) Descriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{3}
}

func (m *Health) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Health.Unmarshal(m, b)
}
func (m *Health) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Health.Marshal(b, m, deterministic)
}
func (m *Health) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Health.Merge(m, src)
}
func (m *Health) XXX_Size() int {
	return xxx_messageInfo_Health.Size(m)
}
func (m *Health) XXX_DiscardUnknown() {
	xxx_messageInfo_Health.DiscardUnknown(m)
}

var xxx_messageInfo_Health proto.InternalMessageInfo

func (m *Health) GetStatus() Health_Status {
	if m != nil {
		return m.Status
	}
	return Health_ALIVE
}

func init() {
	proto.RegisterEnum("core.GerlMsg_Type", GerlMsg_Type_name, GerlMsg_Type_value)
	proto.RegisterEnum("core.Message_Type", Message_Type_name, Message_Type_value)
	proto.RegisterEnum("core.Message_SubType", Message_SubType_name, Message_SubType_value)
	proto.RegisterEnum("core.Health_Status", Health_Status_name, Health_Status_value)
	proto.RegisterType((*Empty)(nil), "core.Empty")
	proto.RegisterType((*GerlMsg)(nil), "core.GerlMsg")
	proto.RegisterType((*Message)(nil), "core.Message")
	proto.RegisterType((*Health)(nil), "core.Health")
}

func init() { proto.RegisterFile("grpc.proto", fileDescriptor_bedfbfc9b54e5600) }

var fileDescriptor_bedfbfc9b54e5600 = []byte{
	// 405 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xdf, 0x8a, 0xd3, 0x40,
	0x14, 0xc6, 0x33, 0x49, 0x36, 0x69, 0x4f, 0xbb, 0x32, 0x9c, 0x45, 0x29, 0xbd, 0x31, 0x8c, 0x8b,
	0x14, 0x85, 0x08, 0xf5, 0x09, 0x96, 0x12, 0x6a, 0xd7, 0x76, 0x5b, 0x26, 0x59, 0xc1, 0xcb, 0xb4,
	0x1d, 0xeb, 0x42, 0x6a, 0x42, 0x66, 0x2a, 0xf4, 0xd2, 0x77, 0xf0, 0x61, 0xbd, 0x94, 0xf9, 0xa3,
	0xa4, 0xe2, 0xdd, 0x39, 0xe7, 0xfb, 0xcd, 0xcc, 0xf7, 0xcd, 0x0c, 0xc0, 0xa1, 0x6d, 0x76, 0x69,
	0xd3, 0xd6, 0xaa, 0xc6, 0x70, 0x57, 0xb7, 0x82, 0xc5, 0x70, 0x95, 0x1d, 0x1b, 0x75, 0x66, 0x3f,
	0x09, 0xc4, 0x73, 0xd1, 0x56, 0x2b, 0x79, 0xc0, 0xd7, 0x10, 0xaa, 0x73, 0x23, 0x46, 0x24, 0x21,
	0x93, 0x67, 0x53, 0x4c, 0x35, 0x99, 0x3a, 0x31, 0x2d, 0xce, 0x8d, 0xe0, 0x46, 0xc7, 0x31, 0xf4,
	0xbe, 0xb4, 0xf5, 0xb1, 0xdc, 0xef, 0xdb, 0x91, 0x9f, 0x90, 0x49, 0x9f, 0xff, 0xed, 0xf1, 0x25,
	0x04, 0x47, 0x79, 0x18, 0x05, 0x09, 0x99, 0x0c, 0xa6, 0xd7, 0x76, 0x8b, 0x95, 0x90, 0xb2, 0x3c,
	0x08, 0xae, 0x15, 0x76, 0x0b, 0xa1, 0xde, 0x0a, 0x7b, 0x10, 0xce, 0xee, 0x96, 0x4b, 0xea, 0xd9,
	0x2a, 0x2f, 0x28, 0xd1, 0xd5, 0x86, 0xaf, 0x67, 0xd4, 0x67, 0xbf, 0x08, 0xc4, 0x6e, 0xd9, 0xff,
	0x6d, 0x39, 0xb1, 0x6b, 0xeb, 0x1d, 0xc4, 0xf2, 0xb4, 0x35, 0xa8, 0x6f, 0xd0, 0xe7, 0x97, 0x68,
	0x7e, 0xda, 0x1a, 0xfa, 0x0f, 0x85, 0x09, 0x0c, 0xf6, 0x42, 0xee, 0xda, 0xa7, 0x46, 0x3d, 0xd5,
	0xdf, 0x8c, 0xe7, 0x3e, 0xef, 0x8e, 0xf0, 0x05, 0x44, 0xdf, 0xcb, 0xea, 0x24, 0xe4, 0x28, 0x4c,
	0x82, 0x49, 0x9f, 0xbb, 0x8e, 0xbd, 0x71, 0x21, 0x00, 0xa2, 0x7c, 0xb1, 0xda, 0x2c, 0x33, 0x1b,
	0x23, 0xff, 0xfc, 0x30, 0xa3, 0x04, 0x87, 0xd0, 0xe3, 0xd9, 0x7c, 0x91, 0x17, 0x19, 0xa7, 0x3e,
	0x4b, 0x21, 0x76, 0x27, 0x63, 0x0c, 0xc1, 0x3c, 0x2b, 0xa8, 0xa7, 0x8b, 0x3c, 0xd3, 0x89, 0x63,
	0x08, 0x36, 0x8f, 0x05, 0xf5, 0xf5, 0xea, 0xfb, 0xf5, 0xe2, 0x81, 0x06, 0xec, 0x1e, 0xa2, 0x0f,
	0xa2, 0xac, 0xd4, 0x57, 0x7c, 0x0b, 0x91, 0x54, 0xa5, 0x3a, 0x49, 0x17, 0xfd, 0xc6, 0xe6, 0xb1,
	0x6a, 0x9a, 0x1b, 0x89, 0x3b, 0x84, 0xdd, 0x40, 0x64, 0x27, 0xd8, 0x87, 0xab, 0xbb, 0xe5, 0xe2,
	0x53, 0x46, 0xbd, 0xe9, 0x0f, 0x02, 0x43, 0xf3, 0x80, 0xf6, 0x0a, 0x5a, 0x7d, 0x97, 0xb3, 0xb2,
	0xaa, 0xf0, 0xfa, 0xe2, 0x71, 0xc7, 0x97, 0x2d, 0xf3, 0xf0, 0x56, 0x73, 0x52, 0xfd, 0xcb, 0x0d,
	0x6c, 0x6b, 0xbf, 0x8e, 0x87, 0xaf, 0x20, 0xe4, 0x8f, 0xeb, 0x8f, 0xd8, 0x1d, 0x8f, 0x87, 0x5d,
	0x97, 0xcc, 0xdb, 0x46, 0xe6, 0xdf, 0xbd, 0xff, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x04, 0x75,
	0x59, 0x85, 0x02, 0x00, 0x00,
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
	RUOK(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Health, error)
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

func (c *gerlMessagerClient) RUOK(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Health, error) {
	out := new(Health)
	err := c.cc.Invoke(ctx, "/core.GerlMessager/RUOK", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GerlMessagerServer is the server API for GerlMessager service.
type GerlMessagerServer interface {
	Call(context.Context, *GerlMsg) (*GerlMsg, error)
	Cast(context.Context, *GerlMsg) (*Empty, error)
	RUOK(context.Context, *Empty) (*Health, error)
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

func _GerlMessager_RUOK_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GerlMessagerServer).RUOK(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.GerlMessager/RUOK",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GerlMessagerServer).RUOK(ctx, req.(*Empty))
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
		{
			MethodName: "RUOK",
			Handler:    _GerlMessager_RUOK_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc.proto",
}
