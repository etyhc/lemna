package client

import (
	fmt "fmt"
	"lemna/agent"
	"lemna/agent/arpc"
)

// Client 客户端rpc调用服务
type Client struct {
	stream arpc.Client_ForwardServer //网络流
	uid    uint32                    //客户端id
	cache  map[uint32]agent.STarget  //转发目标缓存
	Value  interface{}               //使用者可以保存任何数据
}

// NewClient 新客户端
//         s 客户端网络流
//        id 客户端uid，客户端唯一
func NewClient(s arpc.Client_ForwardServer, id uint32) *Client {
	return &Client{stream: s, uid: id, cache: make(map[uint32]agent.STarget)}
}

// Error 附加目标信息到错误上
func (c *Client) Error(err interface{}) error {
	return fmt.Errorf("<uid=%d>%s", c.uid, err)
}

// ID 客户端UID
func (c *Client) ID() uint32 {
	return c.uid
}

// Send 向客户端发送转发消息
func (c *Client) Send(msg *arpc.ForwardMsg) error {
	return c.stream.Send(msg)
}

// Recv 接收客户端转发消息
func (c *Client) Recv() (*arpc.ForwardMsg, error) {
	return c.stream.Recv()
}

// Cache 为客户端提供服务的服务器缓存
func (c *Client) Cache() map[uint32]agent.STarget {
	return c.cache
}
