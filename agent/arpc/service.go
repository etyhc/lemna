package arpc

import (
	"lemna/logger"
	"net"

	"google.golang.org/grpc"
)

type ServerIndex struct {
	index map[string]*Server
}

func NewServerIndex() *ServerIndex {
	return &ServerIndex{index: make(map[string]*Server)}
}

func (si *ServerIndex) get(addr string) *Server {
	return si.index[addr]
}

func (si *ServerIndex) put(s *Server) {
	if addr, ok := s.GetPeerAddr(); ok {
		si.index[addr] = s
	}
}
func (si *ServerIndex) remove(s *Server) {
	if addr, ok := s.GetPeerAddr(); ok {
		delete(si.index, addr)
	}
}

// ServerService rpc服务器端封装，用于服务器开发
type ServerService struct {
	Addr      string     //服务器地址
	Typeid    int32      //服务器类型
	Msgcenter *MsgCenter //消息中心
	si        *ServerIndex
}

func (ss *ServerService) Get(addr string) *Server {
	return ss.si.get(addr)
}

// Forward rpc.Forward调用实现,解析转发来的消息
func (ss *ServerService) Forward(stream Server_ForwardServer) error {
	server := NewServer(ss.Typeid, stream, ss.Msgcenter)
	ss.si.put(server)
	defer ss.si.remove(server)
	for {
		in, err := server.stream.Recv()
		if err == nil {
			err = ss.Msgcenter.Handle(in, server)
			if err != nil {
				//忽略错误的消息
				logger.Error(err)
			}
		} else {
			return err
		}
	}
}

// Run 运行rpc服务,阻塞的
func (ss *ServerService) Run() error {
	lis, err := net.Listen("tcp", ss.Addr)
	logger.Debug(lis.Addr())
	ss.si = NewServerIndex()
	if err == nil {
		rpcs := grpc.NewServer()
		RegisterServerServer(rpcs, ss)
		return rpcs.Serve(lis)
	}
	return err
}
