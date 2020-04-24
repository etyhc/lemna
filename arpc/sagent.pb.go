// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sagent.proto

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
	Mid                  uint32   `protobuf:"varint,2,opt,name=mid,proto3" json:"mid,omitempty"`
	Raw                  []byte   `protobuf:"bytes,3,opt,name=raw,proto3" json:"raw,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MulticastMsg) Reset()         { *m = MulticastMsg{} }
func (m *MulticastMsg) String() string { return proto.CompactTextString(m) }
func (*MulticastMsg) ProtoMessage()    {}
func (*MulticastMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_7449bfcd9a41bcb3, []int{0}
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

func (m *MulticastMsg) GetMid() uint32 {
	if m != nil {
		return m.Mid
	}
	return 0
}

func (m *MulticastMsg) GetRaw() []byte {
	if m != nil {
		return m.Raw
	}
	return nil
}

func init() {
	proto.RegisterType((*MulticastMsg)(nil), "arpc.MulticastMsg")
}

func init() { proto.RegisterFile("sagent.proto", fileDescriptor_7449bfcd9a41bcb3) }

var fileDescriptor_7449bfcd9a41bcb3 = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0x4c, 0x4f,
	0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x2c, 0x2a, 0x48, 0x96, 0xe2,
	0x49, 0x2b, 0x4f, 0xc9, 0x2d, 0x4e, 0x87, 0x88, 0x29, 0xf9, 0x70, 0xf1, 0xf8, 0x96, 0xe6, 0x94,
	0x64, 0x26, 0x27, 0x16, 0x97, 0xf8, 0x16, 0xa7, 0x0b, 0x49, 0x70, 0xb1, 0x97, 0x24, 0x16, 0xa5,
	0xa7, 0x96, 0x14, 0x4b, 0x30, 0x2a, 0x30, 0x6b, 0xf0, 0x06, 0xc1, 0xb8, 0x42, 0x02, 0x5c, 0xcc,
	0xb9, 0x99, 0x29, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xbc, 0x41, 0x20, 0x26, 0x48, 0xa4, 0x28, 0xb1,
	0x5c, 0x82, 0x59, 0x81, 0x51, 0x83, 0x27, 0x08, 0xc4, 0x34, 0x9a, 0xc3, 0xc8, 0xc5, 0x16, 0xec,
	0x08, 0xb2, 0x52, 0xc8, 0x98, 0x8b, 0xdd, 0x2d, 0xbf, 0xa8, 0x3c, 0xb1, 0x28, 0x45, 0x48, 0x40,
	0x0f, 0x64, 0xb1, 0x1e, 0x94, 0xeb, 0x5b, 0x9c, 0x2e, 0x85, 0x21, 0xa2, 0xc4, 0xa0, 0xc1, 0x68,
	0xc0, 0x28, 0x64, 0xce, 0xc5, 0x09, 0x77, 0x8d, 0x90, 0x10, 0x44, 0x11, 0xb2, 0xf3, 0x70, 0x6a,
	0x54, 0xe1, 0x62, 0x71, 0x4e, 0xcc, 0xc9, 0x11, 0xe2, 0x81, 0xc8, 0x07, 0x25, 0x96, 0x83, 0x54,
	0xa3, 0xf0, 0x94, 0x18, 0x92, 0xd8, 0xc0, 0x7e, 0x36, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x5c,
	0x61, 0x09, 0xde, 0x17, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SAgentClient is the client API for SAgent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SAgentClient interface {
	//服务器发送、接收客户端的消息
	Forward(ctx context.Context, opts ...grpc.CallOption) (SAgent_ForwardClient, error)
	//服务器多播消息给客户端、不接收客户端消息
	Multicast(ctx context.Context, opts ...grpc.CallOption) (SAgent_MulticastClient, error)
	//服务器与代理交互rpc
	Call(ctx context.Context, in *RawMsg, opts ...grpc.CallOption) (*RawMsg, error)
}

type sAgentClient struct {
	cc *grpc.ClientConn
}

func NewSAgentClient(cc *grpc.ClientConn) SAgentClient {
	return &sAgentClient{cc}
}

func (c *sAgentClient) Forward(ctx context.Context, opts ...grpc.CallOption) (SAgent_ForwardClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SAgent_serviceDesc.Streams[0], "/arpc.SAgent/Forward", opts...)
	if err != nil {
		return nil, err
	}
	x := &sAgentForwardClient{stream}
	return x, nil
}

type SAgent_ForwardClient interface {
	Send(*ForwardMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ClientStream
}

type sAgentForwardClient struct {
	grpc.ClientStream
}

func (x *sAgentForwardClient) Send(m *ForwardMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sAgentForwardClient) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sAgentClient) Multicast(ctx context.Context, opts ...grpc.CallOption) (SAgent_MulticastClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SAgent_serviceDesc.Streams[1], "/arpc.SAgent/Multicast", opts...)
	if err != nil {
		return nil, err
	}
	x := &sAgentMulticastClient{stream}
	return x, nil
}

type SAgent_MulticastClient interface {
	Send(*MulticastMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ClientStream
}

type sAgentMulticastClient struct {
	grpc.ClientStream
}

func (x *sAgentMulticastClient) Send(m *MulticastMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sAgentMulticastClient) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sAgentClient) Call(ctx context.Context, in *RawMsg, opts ...grpc.CallOption) (*RawMsg, error) {
	out := new(RawMsg)
	err := c.cc.Invoke(ctx, "/arpc.SAgent/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SAgentServer is the server API for SAgent service.
type SAgentServer interface {
	//服务器发送、接收客户端的消息
	Forward(SAgent_ForwardServer) error
	//服务器多播消息给客户端、不接收客户端消息
	Multicast(SAgent_MulticastServer) error
	//服务器与代理交互rpc
	Call(context.Context, *RawMsg) (*RawMsg, error)
}

// UnimplementedSAgentServer can be embedded to have forward compatible implementations.
type UnimplementedSAgentServer struct {
}

func (*UnimplementedSAgentServer) Forward(srv SAgent_ForwardServer) error {
	return status.Errorf(codes.Unimplemented, "method Forward not implemented")
}
func (*UnimplementedSAgentServer) Multicast(srv SAgent_MulticastServer) error {
	return status.Errorf(codes.Unimplemented, "method Multicast not implemented")
}
func (*UnimplementedSAgentServer) Call(ctx context.Context, req *RawMsg) (*RawMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}

func RegisterSAgentServer(s *grpc.Server, srv SAgentServer) {
	s.RegisterService(&_SAgent_serviceDesc, srv)
}

func _SAgent_Forward_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SAgentServer).Forward(&sAgentForwardServer{stream})
}

type SAgent_ForwardServer interface {
	Send(*ForwardMsg) error
	Recv() (*ForwardMsg, error)
	grpc.ServerStream
}

type sAgentForwardServer struct {
	grpc.ServerStream
}

func (x *sAgentForwardServer) Send(m *ForwardMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sAgentForwardServer) Recv() (*ForwardMsg, error) {
	m := new(ForwardMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SAgent_Multicast_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SAgentServer).Multicast(&sAgentMulticastServer{stream})
}

type SAgent_MulticastServer interface {
	Send(*ForwardMsg) error
	Recv() (*MulticastMsg, error)
	grpc.ServerStream
}

type sAgentMulticastServer struct {
	grpc.ServerStream
}

func (x *sAgentMulticastServer) Send(m *ForwardMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sAgentMulticastServer) Recv() (*MulticastMsg, error) {
	m := new(MulticastMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SAgent_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SAgentServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arpc.SAgent/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SAgentServer).Call(ctx, req.(*RawMsg))
	}
	return interceptor(ctx, in, info, handler)
}

var _SAgent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "arpc.SAgent",
	HandlerType: (*SAgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _SAgent_Call_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Forward",
			Handler:       _SAgent_Forward_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Multicast",
			Handler:       _SAgent_Multicast_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "sagent.proto",
}