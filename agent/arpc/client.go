package arpc

import (
	"context"
	fmt "fmt"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
)

// Client rpc客户端服务封装,实现了ServerPeer
type Client struct {
	addr   string //服务器地址
	typeid int32  //服务器类型
	stream Server_ForwardClient
}

func NewClient(addr string, TypeID int32) *Client {
	return &Client{addr: addr, typeid: TypeID}
}

func (c *Client) Forward(target int32, msg proto.Message) error {
	fmsg, err := WrapFMNoCheck(target, msg)
	if err != nil {
		return err
	}
	return c.stream.Send(fmsg)
}

func (c *Client) Recv() (*BroadcastMsg, error) {
	return c.stream.Recv()
}

func (c *Client) Send(fmsg *ForwardMsg) error {
	return c.stream.Send(fmsg)
}

// TypeID 服务器类型
func (c *Client) TypeID() int32 {
	return c.typeid
}

func (c *Client) Error(err interface{}) error {
	return fmt.Errorf("<%s> %s", c.addr, err)
}

// Init 初始化
func (c *Client) Init() error {
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err == nil {
		sc := NewServerClient(conn)
		c.stream, err = sc.Forward(context.Background())
	}
	return err
}
