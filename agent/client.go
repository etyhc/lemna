package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"lemna/logger"
)

//ClientPool 客户端池接口，供ServerPool使用
type ClientPool interface {
	GetClient(int32, *Server) *Client
	SetServerPool(ServerPool)
}

// Client 代理客户端
type Client struct {
	stream rpc.Client_ForwardServer //网络流
	id     int32                    //客户端id
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

//  Run 运行客户端转发功能，循环等待客户端消息并转发给服务器
// pool 转发服务器池
//      等待消息错误，返回
//      在自己的缓存未找到转发服务器，再从转发服务器池寻找转发服务器，无视无转发服务器错误
//      转发失败清除自己缓存的转发服务器
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
		err = s.rpcc.Send(fmsg)
		//转发失败
		if err != nil {
			logger.Error(s.Error(err))
			delete(c.cache, s.rpcc.TypeID())
			continue
		}
	}
}
