// Code generated by protoc-gen-go. DO NOT EDIT.
// source: agent.proto

package rpc

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

type LoginMsg struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginMsg) Reset()         { *m = LoginMsg{} }
func (m *LoginMsg) String() string { return proto.CompactTextString(m) }
func (*LoginMsg) ProtoMessage()    {}
func (*LoginMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_251feda88b728369, []int{0}
}
func (m *LoginMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginMsg.Unmarshal(m, b)
}
func (m *LoginMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginMsg.Marshal(b, m, deterministic)
}
func (dst *LoginMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginMsg.Merge(dst, src)
}
func (m *LoginMsg) XXX_Size() int {
	return xxx_messageInfo_LoginMsg.Size(m)
}
func (m *LoginMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginMsg.DiscardUnknown(m)
}

var xxx_messageInfo_LoginMsg proto.InternalMessageInfo

func (m *LoginMsg) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type RawMsg struct {
	Type                 int32    `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Raw                  []byte   `protobuf:"bytes,2,opt,name=raw,proto3" json:"raw,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawMsg) Reset()         { *m = RawMsg{} }
func (m *RawMsg) String() string { return proto.CompactTextString(m) }
func (*RawMsg) ProtoMessage()    {}
func (*RawMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_251feda88b728369, []int{1}
}
func (m *RawMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawMsg.Unmarshal(m, b)
}
func (m *RawMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawMsg.Marshal(b, m, deterministic)
}
func (dst *RawMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawMsg.Merge(dst, src)
}
func (m *RawMsg) XXX_Size() int {
	return xxx_messageInfo_RawMsg.Size(m)
}
func (m *RawMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_RawMsg.DiscardUnknown(m)
}

var xxx_messageInfo_RawMsg proto.InternalMessageInfo

func (m *RawMsg) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *RawMsg) GetRaw() []byte {
	if m != nil {
		return m.Raw
	}
	return nil
}

type BroadcastMsg struct {
	Targets              []int32  `protobuf:"varint,1,rep,packed,name=targets,proto3" json:"targets,omitempty"`
	Msg                  *RawMsg  `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BroadcastMsg) Reset()         { *m = BroadcastMsg{} }
func (m *BroadcastMsg) String() string { return proto.CompactTextString(m) }
func (*BroadcastMsg) ProtoMessage()    {}
func (*BroadcastMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_251feda88b728369, []int{2}
}
func (m *BroadcastMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BroadcastMsg.Unmarshal(m, b)
}
func (m *BroadcastMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BroadcastMsg.Marshal(b, m, deterministic)
}
func (dst *BroadcastMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BroadcastMsg.Merge(dst, src)
}
func (m *BroadcastMsg) XXX_Size() int {
	return xxx_messageInfo_BroadcastMsg.Size(m)
}
func (m *BroadcastMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_BroadcastMsg.DiscardUnknown(m)
}

var xxx_messageInfo_BroadcastMsg proto.InternalMessageInfo

func (m *BroadcastMsg) GetTargets() []int32 {
	if m != nil {
		return m.Targets
	}
	return nil
}

func (m *BroadcastMsg) GetMsg() *RawMsg {
	if m != nil {
		return m.Msg
	}
	return nil
}

type ForwardMsg struct {
	Target               int32    `protobuf:"varint,1,opt,name=target,proto3" json:"target,omitempty"`
	Msg                  *RawMsg  `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForwardMsg) Reset()         { *m = ForwardMsg{} }
func (m *ForwardMsg) String() string { return proto.CompactTextString(m) }
func (*ForwardMsg) ProtoMessage()    {}
func (*ForwardMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_agent_251feda88b728369, []int{3}
}
func (m *ForwardMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardMsg.Unmarshal(m, b)
}
func (m *ForwardMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardMsg.Marshal(b, m, deterministic)
}
func (dst *ForwardMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardMsg.Merge(dst, src)
}
func (m *ForwardMsg) XXX_Size() int {
	return xxx_messageInfo_ForwardMsg.Size(m)
}
func (m *ForwardMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardMsg proto.InternalMessageInfo

func (m *ForwardMsg) GetTarget() int32 {
	if m != nil {
		return m.Target
	}
	return 0
}

func (m *ForwardMsg) GetMsg() *RawMsg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func init() {
	proto.RegisterType((*LoginMsg)(nil), "rpc.LoginMsg")
	proto.RegisterType((*RawMsg)(nil), "rpc.RawMsg")
	proto.RegisterType((*BroadcastMsg)(nil), "rpc.BroadcastMsg")
	proto.RegisterType((*ForwardMsg)(nil), "rpc.ForwardMsg")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ClientClient is the client API for Client service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClientClient interface {
	Login(ctx context.Context, in *LoginMsg, opts ...grpc.CallOption) (*LoginMsg, error)
	Forward(ctx context.Context, opts ...grpc.CallOption) (Client_ForwardClient, error)
}

type clientClient struct {
	cc *grpc.ClientConn
}

func NewClientClient(cc *grpc.ClientConn) ClientClient {
	return &clientClient{cc}
}

func (c *clientClient) Login(ctx context.Context, in *LoginMsg, opts ...grpc.CallOption) (*LoginMsg, error) {
	out := new(LoginMsg)
	err := c.cc.Invoke(ctx, "/rpc.Client/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) Forward(ctx context.Context, opts ...grpc.CallOption) (Client_ForwardClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Client_serviceDesc.Streams[0], "/rpc.Client/Forward", opts...)
	if err != nil {
		return nil, err
	}
	x := &clientForwardClient{stream}
	return x, nil
}

type Client_ForwardClient interface {
	Send(*ForwardMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ClientStream
}

type clientForwardClient struct {
	grpc.ClientStream
}

func (x *clientForwardClient) Send(m *ForwardMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *clientForwardClient) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ClientServer is the server API for Client service.
type ClientServer interface {
	Login(context.Context, *LoginMsg) (*LoginMsg, error)
	Forward(Client_ForwardServer) error
}

func RegisterClientServer(s *grpc.Server, srv ClientServer) {
	s.RegisterService(&_Client_serviceDesc, srv)
}

func _Client_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Client/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).Login(ctx, req.(*LoginMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_Forward_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ClientServer).Forward(&clientForwardServer{stream})
}

type Client_ForwardServer interface {
	Send(*ForwardMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ServerStream
}

type clientForwardServer struct {
	grpc.ServerStream
}

func (x *clientForwardServer) Send(m *ForwardMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *clientForwardServer) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Client_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Client",
	HandlerType: (*ClientServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Client_Login_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Forward",
			Handler:       _Client_Forward_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "agent.proto",
}

// ServerClient is the client API for Server service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServerClient interface {
	Forward(ctx context.Context, opts ...grpc.CallOption) (Server_ForwardClient, error)
}

type serverClient struct {
	cc *grpc.ClientConn
}

func NewServerClient(cc *grpc.ClientConn) ServerClient {
	return &serverClient{cc}
}

func (c *serverClient) Forward(ctx context.Context, opts ...grpc.CallOption) (Server_ForwardClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Server_serviceDesc.Streams[0], "/rpc.Server/Forward", opts...)
	if err != nil {
		return nil, err
	}
	x := &serverForwardClient{stream}
	return x, nil
}

type Server_ForwardClient interface {
	Send(*ForwardMsg) error
	Recv() (*BroadcastMsg, error)
	grpc.ClientStream
}

type serverForwardClient struct {
	grpc.ClientStream
}

func (x *serverForwardClient) Send(m *ForwardMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serverForwardClient) Recv() (*BroadcastMsg, error) {
	m := new(BroadcastMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServerServer is the server API for Server service.
type ServerServer interface {
	Forward(Server_ForwardServer) error
}

func RegisterServerServer(s *grpc.Server, srv ServerServer) {
	s.RegisterService(&_Server_serviceDesc, srv)
}

func _Server_Forward_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServerServer).Forward(&serverForwardServer{stream})
}

type Server_ForwardServer interface {
	Send(*BroadcastMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ServerStream
}

type serverForwardServer struct {
	grpc.ServerStream
}

func (x *serverForwardServer) Send(m *BroadcastMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serverForwardServer) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Server_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Server",
	HandlerType: (*ServerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Forward",
			Handler:       _Server_Forward_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "agent.proto",
}

func init() { proto.RegisterFile("agent.proto", fileDescriptor_agent_251feda88b728369) }

var fileDescriptor_agent_251feda88b728369 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0x4f, 0x4b, 0xc4, 0x30,
	0x10, 0xc5, 0x37, 0xd6, 0x66, 0x75, 0xba, 0xa2, 0x0e, 0x22, 0x65, 0x41, 0x28, 0xb9, 0xd8, 0x53,
	0xd1, 0xee, 0xd9, 0x8b, 0x0b, 0x7a, 0xd1, 0x4b, 0xfc, 0x04, 0xb1, 0x0d, 0x61, 0x51, 0x9b, 0x32,
	0x0d, 0x16, 0xbf, 0xbd, 0x24, 0x69, 0xf1, 0xcf, 0x41, 0x6f, 0xf3, 0xe6, 0xe5, 0x37, 0x2f, 0x93,
	0x40, 0xa6, 0x8c, 0xee, 0x5c, 0xd5, 0x93, 0x75, 0x16, 0x13, 0xea, 0x1b, 0x51, 0xc0, 0xc1, 0x83,
	0x35, 0xbb, 0xee, 0x71, 0x30, 0x78, 0x06, 0xa9, 0xb3, 0x2f, 0xba, 0xcb, 0x59, 0xc1, 0xca, 0x43,
	0x19, 0x85, 0xa8, 0x80, 0x4b, 0x35, 0x7a, 0x1f, 0x61, 0xdf, 0x7d, 0xf4, 0x3a, 0xd8, 0xa9, 0x0c,
	0x35, 0x9e, 0x40, 0x42, 0x6a, 0xcc, 0xf7, 0x0a, 0x56, 0xae, 0xa4, 0x2f, 0xc5, 0x3d, 0xac, 0x6e,
	0xc9, 0xaa, 0xb6, 0x51, 0x83, 0xf3, 0x54, 0x0e, 0x4b, 0xa7, 0xc8, 0x68, 0x37, 0xe4, 0xac, 0x48,
	0xca, 0x54, 0xce, 0x12, 0x2f, 0x20, 0x79, 0x1b, 0x4c, 0x60, 0xb3, 0x3a, 0xab, 0xa8, 0x6f, 0xaa,
	0x98, 0x24, 0x7d, 0x5f, 0x6c, 0x01, 0xee, 0x2c, 0x8d, 0x8a, 0x5a, 0x3f, 0xe6, 0x1c, 0x78, 0xe4,
	0xa6, 0xf8, 0x49, 0xfd, 0x33, 0xa4, 0x6e, 0x81, 0x6f, 0x5f, 0x77, 0xba, 0x73, 0x78, 0x09, 0x69,
	0xd8, 0x14, 0x8f, 0xc2, 0xa1, 0x79, 0xeb, 0xf5, 0x4f, 0x29, 0x16, 0x78, 0x0d, 0xcb, 0x29, 0x17,
	0x8f, 0x83, 0xf7, 0x75, 0x8b, 0xf5, 0xef, 0x86, 0x58, 0x94, 0xec, 0x8a, 0xd5, 0x37, 0xc0, 0x9f,
	0x34, 0xbd, 0x6b, 0xc2, 0xcd, 0x1f, 0xf0, 0x69, 0x68, 0x7c, 0x7f, 0x9c, 0x88, 0x3f, 0xf3, 0xf0,
	0x21, 0x9b, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x68, 0x9c, 0xca, 0x65, 0x9f, 0x01, 0x00, 0x00,
}
