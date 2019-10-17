// Code generated by protoc-gen-go. DO NOT EDIT.
// source: clientrpc.proto

package arpc

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
	return fileDescriptor_b1249944c579c2db, []int{0}
}

func (m *LoginMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginMsg.Unmarshal(m, b)
}
func (m *LoginMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginMsg.Marshal(b, m, deterministic)
}
func (m *LoginMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginMsg.Merge(m, src)
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

func init() {
	proto.RegisterType((*LoginMsg)(nil), "arpc.LoginMsg")
}

func init() { proto.RegisterFile("clientrpc.proto", fileDescriptor_b1249944c579c2db) }

var fileDescriptor_b1249944c579c2db = []byte{
	// 146 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0xce, 0xc9, 0x4c,
	0xcd, 0x2b, 0x29, 0x2a, 0x48, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x2c, 0x2a,
	0x48, 0x96, 0xe2, 0x49, 0x2b, 0x4f, 0xc9, 0x2d, 0x4e, 0x87, 0x88, 0x29, 0x29, 0x70, 0x71, 0xf8,
	0xe4, 0xa7, 0x67, 0xe6, 0xf9, 0x16, 0xa7, 0x0b, 0x89, 0x70, 0xb1, 0x96, 0xe4, 0x67, 0xa7, 0xe6,
	0x49, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x41, 0x38, 0x46, 0x69, 0x5c, 0x2c, 0xce, 0x45, 0x05,
	0xc9, 0x42, 0x9a, 0x5c, 0xac, 0x60, 0x95, 0x42, 0x7c, 0x7a, 0x20, 0x73, 0xf4, 0x60, 0xda, 0xa4,
	0xd0, 0xf8, 0x4a, 0x0c, 0x42, 0xc6, 0x5c, 0xec, 0x6e, 0xf9, 0x45, 0xe5, 0x89, 0x45, 0x29, 0x42,
	0x02, 0x10, 0x49, 0x28, 0x17, 0xa4, 0x1c, 0x43, 0x44, 0x89, 0x41, 0x83, 0xd1, 0x80, 0x31, 0x89,
	0x0d, 0xec, 0x20, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc9, 0xc7, 0x69, 0xe1, 0xb7, 0x00,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CrpcClient is the client API for Crpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CrpcClient interface {
	Login(ctx context.Context, in *LoginMsg, opts ...grpc.CallOption) (*LoginMsg, error)
	Forward(ctx context.Context, opts ...grpc.CallOption) (Crpc_ForwardClient, error)
}

type crpcClient struct {
	cc *grpc.ClientConn
}

func NewCrpcClient(cc *grpc.ClientConn) CrpcClient {
	return &crpcClient{cc}
}

func (c *crpcClient) Login(ctx context.Context, in *LoginMsg, opts ...grpc.CallOption) (*LoginMsg, error) {
	out := new(LoginMsg)
	err := c.cc.Invoke(ctx, "/arpc.Crpc/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crpcClient) Forward(ctx context.Context, opts ...grpc.CallOption) (Crpc_ForwardClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Crpc_serviceDesc.Streams[0], "/arpc.Crpc/Forward", opts...)
	if err != nil {
		return nil, err
	}
	x := &crpcForwardClient{stream}
	return x, nil
}

type Crpc_ForwardClient interface {
	Send(*ForwardMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ClientStream
}

type crpcForwardClient struct {
	grpc.ClientStream
}

func (x *crpcForwardClient) Send(m *ForwardMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *crpcForwardClient) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CrpcServer is the server API for Crpc service.
type CrpcServer interface {
	Login(context.Context, *LoginMsg) (*LoginMsg, error)
	Forward(Crpc_ForwardServer) error
}

// UnimplementedCrpcServer can be embedded to have forward compatible implementations.
type UnimplementedCrpcServer struct {
}

func (*UnimplementedCrpcServer) Login(ctx context.Context, req *LoginMsg) (*LoginMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (*UnimplementedCrpcServer) Forward(srv Crpc_ForwardServer) error {
	return status.Errorf(codes.Unimplemented, "method Forward not implemented")
}

func RegisterCrpcServer(s *grpc.Server, srv CrpcServer) {
	s.RegisterService(&_Crpc_serviceDesc, srv)
}

func _Crpc_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrpcServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arpc.Crpc/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrpcServer).Login(ctx, req.(*LoginMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Crpc_Forward_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CrpcServer).Forward(&crpcForwardServer{stream})
}

type Crpc_ForwardServer interface {
	Send(*ForwardMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ServerStream
}

type crpcForwardServer struct {
	grpc.ServerStream
}

func (x *crpcForwardServer) Send(m *ForwardMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *crpcForwardServer) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Crpc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "arpc.Crpc",
	HandlerType: (*CrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Crpc_Login_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Forward",
			Handler:       _Crpc_Forward_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "clientrpc.proto",
}
