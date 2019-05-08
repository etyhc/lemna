package client

import (
	fmt "fmt"
	"lemna/agent"
	"lemna/agent/arpc"
	"lemna/logger"
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

// Run 运行客户端转发功能，循环等待客户端消息并转发给服务器
//       等待消息错误，返回
//       在自己的缓存未找到转发服务器，再从转发服务器池寻找转发服务器，无视无转发服务器错误
//       转发失败清除自己缓存的转发服务器
//  pool 转发服务器池
func (c *Client) Run(pool agent.TargetPool) error {
	for {
		fmsg, err := c.stream.Recv()
		if err != nil {
			return c.Error(err)
		}

		s, ok := c.cache[fmsg.Target]
		if !ok {
			s = pool.GetTarget(fmsg.Target).(agent.STarget)
			if s == nil {
				logger.Errorf("not find server<%d>", fmsg.Target)
				continue
			}
		}

		//转发指令
		fmsg.Target = c.uid
		err = s.Send(fmsg)
		//转发失败
		if err != nil {
			logger.Error(s.Error(err))
			delete(c.cache, s.ID())
			continue
		}
	}
}
