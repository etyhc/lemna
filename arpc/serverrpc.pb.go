// Code generated by protoc-gen-go. DO NOT EDIT.
// source: serverrpc.proto

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

//服务器广播消息，代理服务器会将此消息从新封装成Forward消息转给客户端
type MulticastMsg struct {
	Targets              []uint32 `protobuf:"varint,1,rep,packed,name=targets,proto3" json:"targets,omitempty"`
	Msg                  *RawMsg  `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MulticastMsg) Reset()         { *m = MulticastMsg{} }
func (m *MulticastMsg) String() string { return proto.CompactTextString(m) }
func (*MulticastMsg) ProtoMessage()    {}
func (*MulticastMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd512e99445e691d, []int{0}
}

func (m *MulticastMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MulticastMsg.Unmarshal(m, b)
}
func (m *MulticastMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MulticastMsg.Marshal(b, m, deterministic)
}
func (m *MulticastMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MulticastMsg.Merge(m, src)
}
func (m *MulticastMsg) XXX_Size() int {
	return xxx_messageInfo_MulticastMsg.Size(m)
}
func (m *MulticastMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_MulticastMsg.DiscardUnknown(m)
}

var xxx_messageInfo_MulticastMsg proto.InternalMessageInfo

func (m *MulticastMsg) GetTargets() []uint32 {
	if m != nil {
		return m.Targets
	}
	return nil
}

func (m *MulticastMsg) GetMsg() *RawMsg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func init() {
	proto.RegisterType((*MulticastMsg)(nil), "arpc.MulticastMsg")
}

func init() { proto.RegisterFile("serverrpc.proto", fileDescriptor_dd512e99445e691d) }

var fileDescriptor_dd512e99445e691d = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0x2a, 0x2a, 0x48, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x2c, 0x2a,
	0x48, 0x96, 0xe2, 0x49, 0x2b, 0x4f, 0xc9, 0x2d, 0x4e, 0x87, 0x88, 0x29, 0x79, 0x70, 0xf1, 0xf8,
	0x96, 0xe6, 0x94, 0x64, 0x26, 0x27, 0x16, 0x97, 0xf8, 0x16, 0xa7, 0x0b, 0x49, 0x70, 0xb1, 0x97,
	0x24, 0x16, 0xa5, 0xa7, 0x96, 0x14, 0x4b, 0x30, 0x2a, 0x30, 0x6b, 0xf0, 0x06, 0xc1, 0xb8, 0x42,
	0x72, 0x5c, 0xcc, 0xb9, 0xc5, 0xe9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0x3c, 0x7a, 0x20,
	0xb3, 0xf4, 0x82, 0x12, 0xcb, 0x7d, 0x8b, 0xd3, 0x83, 0x40, 0x12, 0x46, 0xf3, 0x19, 0xb9, 0x58,
	0x82, 0x8b, 0x0a, 0x92, 0x85, 0x8c, 0xb9, 0xd8, 0xdd, 0xf2, 0x8b, 0xca, 0x13, 0x8b, 0x52, 0x84,
	0x04, 0x20, 0xca, 0xa0, 0x5c, 0xdf, 0xe2, 0x74, 0x29, 0x0c, 0x11, 0x25, 0x06, 0x0d, 0x46, 0x03,
	0x46, 0x21, 0x73, 0x2e, 0x4e, 0xb8, 0x3b, 0x84, 0x84, 0x20, 0x8a, 0x90, 0x1d, 0x86, 0x53, 0xa3,
	0x26, 0x17, 0xab, 0x7f, 0x49, 0x46, 0x6a, 0x91, 0x10, 0x8a, 0x93, 0xa4, 0x50, 0x78, 0x10, 0xa5,
	0x49, 0x6c, 0x60, 0x2f, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x07, 0xf6, 0x03, 0xaf, 0x19,
	0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SrpcClient is the client API for Srpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SrpcClient interface {
	//服务器发送、接收客户端的消息
	Forward(ctx context.Context, opts ...grpc.CallOption) (Srpc_ForwardClient, error)
	//服务器多播消息给客户端、不接收客户端消息
	Multicast(ctx context.Context, opts ...grpc.CallOption) (Srpc_MulticastClient, error)
	//服务器与代理交互rpc
	Other(ctx context.Context, opts ...grpc.CallOption) (Srpc_OtherClient, error)
}

type srpcClient struct {
	cc *grpc.ClientConn
}

func NewSrpcClient(cc *grpc.ClientConn) SrpcClient {
	return &srpcClient{cc}
}

func (c *srpcClient) Forward(ctx context.Context, opts ...grpc.CallOption) (Srpc_ForwardClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Srpc_serviceDesc.Streams[0], "/arpc.Srpc/Forward", opts...)
	if err != nil {
		return nil, err
	}
	x := &srpcForwardClient{stream}
	return x, nil
}

type Srpc_ForwardClient interface {
	Send(*ForwardMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ClientStream
}

type srpcForwardClient struct {
	grpc.ClientStream
}

func (x *srpcForwardClient) Send(m *ForwardMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *srpcForwardClient) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *srpcClient) Multicast(ctx context.Context, opts ...grpc.CallOption) (Srpc_MulticastClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Srpc_serviceDesc.Streams[1], "/arpc.Srpc/Multicast", opts...)
	if err != nil {
		return nil, err
	}
	x := &srpcMulticastClient{stream}
	return x, nil
}

type Srpc_MulticastClient interface {
	Send(*MulticastMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ClientStream
}

type srpcMulticastClient struct {
	grpc.ClientStream
}

func (x *srpcMulticastClient) Send(m *MulticastMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *srpcMulticastClient) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *srpcClient) Other(ctx context.Context, opts ...grpc.CallOption) (Srpc_OtherClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Srpc_serviceDesc.Streams[2], "/arpc.Srpc/Other", opts...)
	if err != nil {
		return nil, err
	}
	x := &srpcOtherClient{stream}
	return x, nil
}

type Srpc_OtherClient interface {
	Send(*RawMsg) error
	Recv() (*RawMsg, error)
	grpc.ClientStream
}

type srpcOtherClient struct {
	grpc.ClientStream
}

func (x *srpcOtherClient) Send(m *RawMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *srpcOtherClient) Recv() (*RawMsg, error) {
	m := new(RawMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SrpcServer is the server API for Srpc service.
type SrpcServer interface {
	//服务器发送、接收客户端的消息
	Forward(Srpc_ForwardServer) error
	//服务器多播消息给客户端、不接收客户端消息
	Multicast(Srpc_MulticastServer) error
	//服务器与代理交互rpc
	Other(Srpc_OtherServer) error
}

// UnimplementedSrpcServer can be embedded to have forward compatible implementations.
type UnimplementedSrpcServer struct {
}

func (*UnimplementedSrpcServer) Forward(srv Srpc_ForwardServer) error {
	return status.Errorf(codes.Unimplemented, "method Forward not implemented")
}
func (*UnimplementedSrpcServer) Multicast(srv Srpc_MulticastServer) error {
	return status.Errorf(codes.Unimplemented, "method Multicast not implemented")
}
func (*UnimplementedSrpcServer) Other(srv Srpc_OtherServer) error {
	return status.Errorf(codes.Unimplemented, "method Other not implemented")
}

func RegisterSrpcServer(s *grpc.Server, srv SrpcServer) {
	s.RegisterService(&_Srpc_serviceDesc, srv)
}

func _Srpc_Forward_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SrpcServer).Forward(&srpcForwardServer{stream})
}

type Srpc_ForwardServer interface {
	Send(*ForwardMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ServerStream
}

type srpcForwardServer struct {
	grpc.ServerStream
}

func (x *srpcForwardServer) Send(m *ForwardMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *srpcForwardServer) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Srpc_Multicast_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SrpcServer).Multicast(&srpcMulticastServer{stream})
}

type Srpc_MulticastServer interface {
	Send(*ForwardMsg) error
	Recv() (*MulticastMsg, error)
	grpc.ServerStream
}

type srpcMulticastServer struct {
	grpc.ServerStream
}

func (x *srpcMulticastServer) Send(m *ForwardMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *srpcMulticastServer) Recv() (*MulticastMsg, error) {
	m := new(MulticastMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Srpc_Other_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SrpcServer).Other(&srpcOtherServer{stream})
}

type Srpc_OtherServer interface {
	Send(*RawMsg) error
	Recv() (*RawMsg, error)
	grpc.ServerStream
}

type srpcOtherServer struct {
	grpc.ServerStream
}

func (x *srpcOtherServer) Send(m *RawMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *srpcOtherServer) Recv() (*RawMsg, error) {
	m := new(RawMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Srpc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "arpc.Srpc",
	HandlerType: (*SrpcServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Forward",
			Handler:       _Srpc_Forward_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Multicast",
			Handler:       _Srpc_Multicast_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Other",
			Handler:       _Srpc_Other_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "serverrpc.proto",
}