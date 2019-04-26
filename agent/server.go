package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"lemna/config"
	"lemna/logger"
)

type ServerPool interface {
	GetServer(int32, *Client) *Server
	SetClientPool(cp ClientPool)
}

type Server struct {
	stream rpc.Server_ForwardClient //目标网络流
	typeid int32
	Info   *config.ServerInfo
	Round  int32
}

func NewServer(stream rpc.Server_ForwardClient, tid int32, info *config.ServerInfo) *Server {
	return &Server{stream: stream, typeid: tid, Info: info}
}
func (s *Server) Error(err interface{}) error {
	return fmt.Errorf("<id=%d>%s", s.typeid, err)
}

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
