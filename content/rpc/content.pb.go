// Code generated by protoc-gen-go. DO NOT EDIT.
// source: content.proto

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

type ContentMsg struct {
	Info                 string   `protobuf:"bytes,1,opt,name=Info,proto3" json:"Info,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContentMsg) Reset()         { *m = ContentMsg{} }
func (m *ContentMsg) String() string { return proto.CompactTextString(m) }
func (*ContentMsg) ProtoMessage()    {}
func (*ContentMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_9c9c57ff2cc4a0f3, []int{0}
}
func (m *ContentMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContentMsg.Unmarshal(m, b)
}
func (m *ContentMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContentMsg.Marshal(b, m, deterministic)
}
func (dst *ContentMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContentMsg.Merge(dst, src)
}
func (m *ContentMsg) XXX_Size() int {
	return xxx_messageInfo_ContentMsg.Size(m)
}
func (m *ContentMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ContentMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ContentMsg proto.InternalMessageInfo

func (m *ContentMsg) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func (m *ContentMsg) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*ContentMsg)(nil), "rpc.ContentMsg")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChannelClient is the client API for Channel service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChannelClient interface {
	Publish(ctx context.Context, in *ContentMsg, opts ...grpc.CallOption) (*ContentMsg, error)
	Subscribe(ctx context.Context, in *ContentMsg, opts ...grpc.CallOption) (Channel_SubscribeClient, error)
}

type channelClient struct {
	cc *grpc.ClientConn
}

func NewChannelClient(cc *grpc.ClientConn) ChannelClient {
	return &channelClient{cc}
}

func (c *channelClient) Publish(ctx context.Context, in *ContentMsg, opts ...grpc.CallOption) (*ContentMsg, error) {
	out := new(ContentMsg)
	err := c.cc.Invoke(ctx, "/rpc.Channel/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelClient) Subscribe(ctx context.Context, in *ContentMsg, opts ...grpc.CallOption) (Channel_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Channel_serviceDesc.Streams[0], "/rpc.Channel/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &channelSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Channel_SubscribeClient interface {
	Recv() (*ContentMsg, error)
	grpc.ClientStream
}

type channelSubscribeClient struct {
	grpc.ClientStream
}

func (x *channelSubscribeClient) Recv() (*ContentMsg, error) {
	m := new(ContentMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChannelServer is the server API for Channel service.
type ChannelServer interface {
	Publish(context.Context, *ContentMsg) (*ContentMsg, error)
	Subscribe(*ContentMsg, Channel_SubscribeServer) error
}

func RegisterChannelServer(s *grpc.Server, srv ChannelServer) {
	s.RegisterService(&_Channel_serviceDesc, srv)
}

func _Channel_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContentMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Channel/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelServer).Publish(ctx, req.(*ContentMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Channel_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ContentMsg)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChannelServer).Subscribe(m, &channelSubscribeServer{stream})
}

type Channel_SubscribeServer interface {
	Send(*ContentMsg) error
	grpc.ServerStream
}

type channelSubscribeServer struct {
	grpc.ServerStream
}

func (x *channelSubscribeServer) Send(m *ContentMsg) error {
	return x.ServerStream.SendMsg(m)
}

var _Channel_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Channel",
	HandlerType: (*ChannelServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _Channel_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Channel_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "content.proto",
}

func init() { proto.RegisterFile("content.proto", fileDescriptor_content_9c9c57ff2cc4a0f3) }

var fileDescriptor_content_9c9c57ff2cc4a0f3 = []byte{
	// 141 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0xce, 0xcf, 0x2b,
	0x49, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x2a, 0x48, 0x56, 0x32,
	0xe1, 0xe2, 0x72, 0x86, 0x88, 0xfa, 0x16, 0xa7, 0x0b, 0x09, 0x71, 0xb1, 0x78, 0xe6, 0xa5, 0xe5,
	0x4b, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x20, 0xb1, 0xbc, 0xc4, 0xdc, 0x54, 0x09,
	0x26, 0x88, 0x18, 0x88, 0x6d, 0x94, 0xcd, 0xc5, 0xee, 0x9c, 0x91, 0x98, 0x97, 0x97, 0x9a, 0x23,
	0xa4, 0xcb, 0xc5, 0x1e, 0x50, 0x9a, 0x94, 0x93, 0x59, 0x9c, 0x21, 0xc4, 0xaf, 0x57, 0x54, 0x90,
	0xac, 0x87, 0x30, 0x4e, 0x0a, 0x5d, 0x40, 0x89, 0x41, 0xc8, 0x90, 0x8b, 0x33, 0xb8, 0x34, 0xa9,
	0x38, 0xb9, 0x28, 0x33, 0x29, 0x95, 0x18, 0x0d, 0x06, 0x8c, 0x49, 0x6c, 0x60, 0xe7, 0x1a, 0x03,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x02, 0x36, 0x99, 0xa6, 0xbf, 0x00, 0x00, 0x00,
}