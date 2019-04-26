package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"lemna/logger"
)

type ClientPool interface {
	GetClient(cid int32, s *Server) *Client
	SetServerPool(sp ServerPool)
}

// Client 代理服务的对象统称目标，目标可以相互转发消息
type Client struct {
	stream rpc.Client_ForwardServer //目标网络流
	id     int32                    //目标id
	cache  map[int32]*Server        //转发目标缓存
	Value  interface{}              //使用者可以保存任何数据
}

// NewClient 新代理目标
//         s 目标的网络流
//        id 目标标识，客户端唯一，服务器可能不唯一
func NewClient(s rpc.Client_ForwardServer, id int32) *Client {
	return &Client{stream: s, id: id, cache: make(map[int32]*Server)}
}

// Error 附加目标信息到错误上
func (c *Client) Error(err interface{}) error {
	return fmt.Errorf("<id=%d>%s", c.id, err)
}

// Run  运行转发功能，循环等待消息并转发
// pool 转发目标池
//      等待消息错误返回
//      在自己的缓存未找到转发目标，再从转发目标池寻找转发目标，无视无转发目标错误
//      将自己缓存到转发目标
//      转发失败清除自己缓存的转发目标
func (c *Client) Run(pool ServerPool) error {
	for {
		fmsg, err := c.stream.Recv()
		if err != nil {
			return c.Error(err)
		}

		s, ok := c.cache[fmsg.Target]
		if !ok {
			s = pool.GetServer(fmsg.Target, c)
			if s == nil {
				logger.Errorf("not find server<%d>", fmsg.Target)
				continue
			}
		}

		//转发指令
		fmsg.Target = c.id
		err = s.stream.Send(fmsg)
		//转发失败
		if err != nil {
			logger.Error(s.Error(err))
			delete(c.cache, s.typeid)
			continue
		}
	}
}
