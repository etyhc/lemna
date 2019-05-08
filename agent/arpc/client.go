package arpc

import (
	"context"
	fmt "fmt"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
)

// Client 服务器rpc客户端
type Client struct {
	addr   string               //服务器地址
	typeid uint32               //服务器类型
	stream Server_ForwardClient //服务器流
	id     uint32               //代理唯一ID，每个Client都一样
}

// NewClient 创建一个服务器rpc客户端
//      addr 服务器地址
//    TypeID 服务器类型
//        id 代理服务器唯一ID
func NewClient(addr string, TypeID uint32, id uint32) *Client {
	return &Client{addr: addr, typeid: TypeID, id: id}
}

// Forward 向服务器发送消息
//         这些消息不是客户端来的，是代理服务器需要发送服务器的消息，应用场景比如：
//         客户端不可用，代理通知为此用户提供服务的服务器，此用户不可用
//  target 客户端ID
//     msg 要发送的消息
func (c *Client) Forward(target uint32, msg proto.Message) error {
	fmsg, err := WrapFMNoCheck(target, msg)
	if err != nil {
		return err
	}
	return c.stream.Send(fmsg)
}

// Recv 接收服务器的广播消息
func (c *Client) Recv() (*BroadcastMsg, error) {
	return c.stream.Recv()
}

// Send 向服务器发送客户端转发消息
func (c *Client) Send(fmsg *ForwardMsg) error {
	return c.stream.Send(fmsg)
}

// ID 代理服务器唯一ID，所有服务器rpc客户端都一样
func (c *Client) ID() uint32 {
	return c.id
}

// TypeID 服务器类型
func (c *Client) TypeID() uint32 {
	return c.typeid
}

// Error 在错误信息上附加服务器信息
func (c *Client) Error(err interface{}) error {
	return fmt.Errorf("<%d:%s> %s", c.typeid, c.addr, err)
}

// GetRequestMetadata 接口实现
//                    目的是将代理服务器唯一ID以Metadata形式发送给服务器，供服务器识别代理服务器
func (c *Client) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"clientid": fmt.Sprint(c.id)}, nil
}

// RequireTransportSecurity 接口实现
func (c *Client) RequireTransportSecurity() bool {
	return false
}

// Init 初始化，连接服务器，并发起rpc调用
func (c *Client) Init() error {
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure(), grpc.WithPerRPCCredentials(c))
	if err == nil {
		sc := NewServerClient(conn)
		c.stream, err = sc.Forward(context.Background())
	}
	return err
}
