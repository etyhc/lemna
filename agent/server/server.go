package server

import (
	"lemna/agent/arpc"
)

//Server 服务器
type Server struct {
	rpcc  *arpc.Client
	Info  *ServerInfo //服务器信息
	Round int32       //服务器被调度次数
}

//NewServer 新服务器
func NewServer(client *arpc.Client, info *ServerInfo) *Server {
	return &Server{rpcc: client, Info: info}
}

//Error 附加服务器信息到错误信息上
func (s *Server) Error(err interface{}) error {
	return s.rpcc.Error(err)
}

func (s *Server) Send(msg *arpc.ForwardMsg) error {
	return s.rpcc.Send(msg)
}

func (s *Server) ID() int32 {
	return s.rpcc.TypeID()
}

func (s *Server) Recv() (*arpc.BroadcastMsg, error) {
	return s.rpcc.Recv()
}
