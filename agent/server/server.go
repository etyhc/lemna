package server

import (
	"lemna/agent/arpc"
)

// Server 代理服务器rpc调用客户端
type Server struct {
	rpcc  *arpc.Client //rpc客户端
	Info  *ServerInfo  //服务器信息
	Round int32        //服务器被调度次数
}

// NewServer 新服务器
//    client rpc客户端
//      info 订阅的服务器信息
func NewServer(client *arpc.Client, info *ServerInfo) *Server {
	return &Server{rpcc: client, Info: info}
}

// Error 附加服务器信息到错误信息上
func (s *Server) Error(err interface{}) error {
	return s.rpcc.Error(err)
}

// Send 发送转发消息给服务器
func (s *Server) Send(msg *arpc.ForwardMsg) error {
	return s.rpcc.Send(msg)
}

// ID 服务器类型ID
func (s *Server) ID() int32 {
	return s.rpcc.TypeID()
}

// Recv 接收服务器的广播消息
func (s *Server) Recv() (*arpc.BroadcastMsg, error) {
	return s.rpcc.Recv()
}
