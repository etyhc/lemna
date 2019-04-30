package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"lemna/logger"
)

//ServerPool 服务器池接口，供ClientPool使用
type ServerPool interface {
	GetServer(int32, *Client) *Server
	SetClientPool(cp ClientPool)
}

//Server 服务器
type Server struct {
	stream rpc.Server_ForwardClient //服务器网络流
	typeid int32                    //服务器类型
	Info   *ServerInfo              //服务器信息
	Round  int32                    //服务器被调度次数
}

//NewServer 新服务器
func NewServer(stream rpc.Server_ForwardClient, tid int32, info *ServerInfo) *Server {
	return &Server{stream: stream, typeid: tid, Info: info}
}

//Error 附加服务器信息到错误信息上
func (s *Server) Error(err interface{}) error {
	return fmt.Errorf("<id=%d>%s", s.typeid, err)
}

//Run 运行服务器,接收服务器消息，转发给客户端
func (s *Server) Run(pool ClientPool) error {
	for {
		bmsg, err := s.stream.Recv()
		if err != nil {
			return s.Error(err)
		}

		for _, cid := range bmsg.Targets {
			c := pool.GetClient(cid, s)
			if c == nil {
				logger.Errorf("not find client<%d>", cid)
				continue
			}

			//转发指令
			err = c.stream.Send(&rpc.ForwardMsg{Target: s.typeid, Msg: bmsg.Msg})
			//转发失败
			if err != nil {
				logger.Error(c.Error(err))
				delete(c.cache, s.typeid)
			} else {
				c.cache[s.typeid] = s
			}
		}
	}
}
