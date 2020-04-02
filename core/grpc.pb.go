// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc.proto

package core

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type GerlMsg_Type int32

const (
	GerlMsg_CALL GerlMsg_Type = 0
	GerlMsg_CAST GerlMsg_Type = 1
	GerlMsg_PROC GerlMsg_Type = 2
	GerlMsg_ERR  GerlMsg_Type = 3
	GerlMsg_TERM GerlMsg_Type = 4
)

var GerlMsg_Type_name = map[int32]string{
	0: "CALL",
	1: "CAST",
	2: "PROC",
	3: "ERR",
	4: "TERM",
}

var GerlMsg_Type_value = map[string]int32{
	"CALL": 0,
	"CAST": 1,
	"PROC": 2,
	"ERR":  3,
	"TERM": 4,
}

func (x GerlMsg_Type) String() string {
	return proto.EnumName(GerlMsg_Type_name, int32(x))
}

func (GerlMsg_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_bedfbfc9b54e5600, []int{1, 0}
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
	RawMsg               []byte   `protobuf:"bytes,1,opt,name=rawMsg,proto3" json:"rawMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
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

func (m *Message) GetRawMsg() []byte {
	if m != nil {
		return m.RawMsg
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
	proto.RegisterEnum("core.Health_Status", Health_Status_name, Health_Status_value)
	proto.RegisterType((*Empty)(nil), "core.Empty")
	proto.RegisterType((*GerlMsg)(nil), "core.GerlMsg")
	proto.RegisterType((*Message)(nil), "core.Message")
	proto.RegisterType((*Health)(nil), "core.Health")
}

func init() { proto.RegisterFile("grpc.proto", fileDescriptor_bedfbfc9b54e5600) }

var fileDescriptor_bedfbfc9b54e5600 = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0x4d, 0x4b, 0xc3, 0x40,
	0x14, 0xcc, 0x36, 0x69, 0xd2, 0xbe, 0xb6, 0xb2, 0xbc, 0x82, 0x94, 0x5e, 0xac, 0xab, 0x48, 0x41,
	0xc8, 0xa1, 0x1e, 0x3c, 0x97, 0x12, 0xfc, 0x6a, 0xa8, 0x6c, 0xab, 0xf7, 0xb5, 0x5d, 0xe3, 0x21,
	0x25, 0x61, 0x77, 0x45, 0x72, 0xf4, 0xd7, 0xf8, 0x37, 0x65, 0xb3, 0x41, 0x5a, 0x6f, 0x6f, 0xde,
	0x0c, 0xb3, 0xf3, 0x66, 0x01, 0x32, 0x55, 0x6e, 0xe3, 0x52, 0x15, 0xa6, 0xc0, 0x60, 0x5b, 0x28,
	0xc9, 0x22, 0x68, 0x27, 0xfb, 0xd2, 0x54, 0xec, 0x87, 0x40, 0x74, 0x27, 0x55, 0x9e, 0xea, 0x0c,
	0xaf, 0x20, 0x30, 0x55, 0x29, 0x47, 0x64, 0x42, 0xa6, 0x27, 0x33, 0x8c, 0xad, 0x32, 0x6e, 0xc8,
	0x78, 0x53, 0x95, 0x92, 0xd7, 0x3c, 0x8e, 0xa1, 0xf3, 0xae, 0x8a, 0xbd, 0xd8, 0xed, 0xd4, 0xa8,
	0x35, 0x21, 0xd3, 0x2e, 0xff, 0xc3, 0x78, 0x06, 0xfe, 0x5e, 0x67, 0x23, 0x7f, 0x42, 0xa6, 0xbd,
	0xd9, 0xc0, 0x59, 0xa4, 0x52, 0x6b, 0x91, 0x49, 0x6e, 0x19, 0x76, 0x0b, 0x81, 0xb5, 0xc2, 0x0e,
	0x04, 0x8b, 0xf9, 0x72, 0x49, 0x3d, 0x37, 0xad, 0x37, 0x94, 0xd8, 0xe9, 0x99, 0xaf, 0x16, 0xb4,
	0x85, 0x11, 0xf8, 0x09, 0xe7, 0xd4, 0xb7, 0xab, 0x4d, 0xc2, 0x53, 0x1a, 0xb0, 0x73, 0x88, 0x1a,
	0x23, 0x3c, 0x85, 0x50, 0x89, 0xaf, 0x54, 0x67, 0x75, 0xd4, 0x3e, 0x6f, 0x10, 0x7b, 0x84, 0xf0,
	0x5e, 0x8a, 0xdc, 0x7c, 0xe0, 0x35, 0x84, 0xda, 0x08, 0xf3, 0xa9, 0x9b, 0x63, 0x86, 0x2e, 0x89,
	0x63, 0xe3, 0x75, 0x4d, 0xf1, 0x46, 0xc2, 0x86, 0x10, 0xba, 0x0d, 0x76, 0xa1, 0x3d, 0x5f, 0x3e,
	0xbc, 0x26, 0xd4, 0x9b, 0x7d, 0x13, 0xe8, 0xd7, 0xb7, 0xbb, 0x37, 0x95, 0x6d, 0x67, 0x21, 0xf2,
	0x1c, 0x07, 0x47, 0xbd, 0x8c, 0x8f, 0x21, 0xf3, 0xf0, 0xd2, 0xea, 0xb4, 0xf9, 0xaf, 0xeb, 0x39,
	0xe8, 0x5a, 0xf7, 0xf0, 0x02, 0x02, 0xfe, 0xb2, 0x7a, 0xc2, 0xc3, 0xf5, 0xb8, 0x7f, 0x98, 0x92,
	0x79, 0x6f, 0x61, 0xfd, 0x65, 0x37, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x5b, 0xcc, 0x43, 0x99,
	0xc0, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// GerlMessagerClient is the client API for GerlMessager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GerlMessagerClient interface {
	Call(ctx context.Context, in *GerlMsg, opts ...grpc.CallOption) (*GerlMsg, error)
	Cast(ctx context.Context, in *GerlMsg, opts ...grpc.CallOption) (*Empty, error)
	RUOK(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Health, error)
}

type gerlMessagerClient struct {
	cc grpc.ClientConnInterface
}

func NewGerlMessagerClient(cc grpc.ClientConnInterface) GerlMessagerClient {
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

// UnimplementedGerlMessagerServer can be embedded to have forward compatible implementations.
type UnimplementedGerlMessagerServer struct {
}

func (*UnimplementedGerlMessagerServer) Call(ctx context.Context, req *GerlMsg) (*GerlMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (*UnimplementedGerlMessagerServer) Cast(ctx context.Context, req *GerlMsg) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cast not implemented")
}
func (*UnimplementedGerlMessagerServer) RUOK(ctx context.Context, req *Empty) (*Health, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RUOK not implemented")
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
