// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

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

type ConfigMsg struct {
	Info                 string   `protobuf:"bytes,1,opt,name=Info,proto3" json:"Info,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigMsg) Reset()         { *m = ConfigMsg{} }
func (m *ConfigMsg) String() string { return proto.CompactTextString(m) }
func (*ConfigMsg) ProtoMessage()    {}
func (*ConfigMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_ac5cd3ed9cf118bc, []int{0}
}
func (m *ConfigMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigMsg.Unmarshal(m, b)
}
func (m *ConfigMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigMsg.Marshal(b, m, deterministic)
}
func (dst *ConfigMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigMsg.Merge(dst, src)
}
func (m *ConfigMsg) XXX_Size() int {
	return xxx_messageInfo_ConfigMsg.Size(m)
}
func (m *ConfigMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigMsg proto.InternalMessageInfo

func (m *ConfigMsg) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func (m *ConfigMsg) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*ConfigMsg)(nil), "rpc.ConfigMsg")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConfigClient is the client API for Config service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConfigClient interface {
	Publish(ctx context.Context, in *ConfigMsg, opts ...grpc.CallOption) (*ConfigMsg, error)
	Subscribe(ctx context.Context, in *ConfigMsg, opts ...grpc.CallOption) (Config_SubscribeClient, error)
}

type configClient struct {
	cc *grpc.ClientConn
}

func NewConfigClient(cc *grpc.ClientConn) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) Publish(ctx context.Context, in *ConfigMsg, opts ...grpc.CallOption) (*ConfigMsg, error) {
	out := new(ConfigMsg)
	err := c.cc.Invoke(ctx, "/rpc.Config/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) Subscribe(ctx context.Context, in *ConfigMsg, opts ...grpc.CallOption) (Config_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Config_serviceDesc.Streams[0], "/rpc.Config/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &configSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Config_SubscribeClient interface {
	Recv() (*ConfigMsg, error)
	grpc.ClientStream
}

type configSubscribeClient struct {
	grpc.ClientStream
}

func (x *configSubscribeClient) Recv() (*ConfigMsg, error) {
	m := new(ConfigMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConfigServer is the server API for Config service.
type ConfigServer interface {
	Publish(context.Context, *ConfigMsg) (*ConfigMsg, error)
	Subscribe(*ConfigMsg, Config_SubscribeServer) error
}

func RegisterConfigServer(s *grpc.Server, srv ConfigServer) {
	s.RegisterService(&_Config_serviceDesc, srv)
}

func _Config_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Config/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).Publish(ctx, req.(*ConfigMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConfigMsg)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConfigServer).Subscribe(m, &configSubscribeServer{stream})
}

type Config_SubscribeServer interface {
	Send(*ConfigMsg) error
	grpc.ServerStream
}

type configSubscribeServer struct {
	grpc.ServerStream
}

func (x *configSubscribeServer) Send(m *ConfigMsg) error {
	return x.ServerStream.SendMsg(m)
}

var _Config_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Config",
	HandlerType: (*ConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _Config_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Config_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "config.proto",
}

func init() { proto.RegisterFile("config.proto", fileDescriptor_config_ac5cd3ed9cf118bc) }

var fileDescriptor_config_ac5cd3ed9cf118bc = []byte{
	// 135 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0x4b,
	0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x2a, 0x48, 0x56, 0x32, 0xe6,
	0xe2, 0x74, 0x06, 0x0b, 0xfa, 0x16, 0xa7, 0x0b, 0x09, 0x71, 0xb1, 0x78, 0xe6, 0xa5, 0xe5, 0x4b,
	0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x20, 0xb1, 0xbc, 0xc4, 0xdc, 0x54, 0x09, 0x26,
	0x88, 0x18, 0x88, 0x6d, 0x94, 0xc6, 0xc5, 0x06, 0xd1, 0x24, 0xa4, 0xcd, 0xc5, 0x1e, 0x50, 0x9a,
	0x94, 0x93, 0x59, 0x9c, 0x21, 0xc4, 0xa7, 0x57, 0x54, 0x90, 0xac, 0x07, 0x37, 0x4c, 0x0a, 0x8d,
	0xaf, 0xc4, 0x20, 0xa4, 0xcf, 0xc5, 0x19, 0x5c, 0x9a, 0x54, 0x9c, 0x5c, 0x94, 0x99, 0x94, 0x4a,
	0x58, 0xb9, 0x01, 0x63, 0x12, 0x1b, 0xd8, 0xa1, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb0,
	0x62, 0xc8, 0x68, 0xb8, 0x00, 0x00, 0x00,
}
