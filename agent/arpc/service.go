// Package arpc 基于grpc的代理服务器的基础定义.
//
//             定义了2个grpc服务Outside和Inside
//             客户端<----------->代理<---------->服务器
//                    Outside rpc      Inside rpc
package arpc

import (
	fmt "fmt"
	"lemna/logger"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServerIndex struct {
	index map[uint32]*Server
}

func NewServerIndex() *ServerIndex {
	return &ServerIndex{index: make(map[uint32]*Server)}
}

func (si *ServerIndex) get(id uint32) *Server {
	return si.index[id]
}

func (si *ServerIndex) put(s *Server) error {
	if _, ok := si.index[s.ID()]; ok {
		return fmt.Errorf("clientid<%d> conflcit", s.ID())
	}
	si.index[s.ID()] = s
	return nil
}

func (si *ServerIndex) remove(s *Server) {
	delete(si.index, s.ID())
}

// ServerService rpc服务器端封装，用于服务器开发
type ServerService struct {
	Addr      string     //服务器地址
	Typeid    int32      //服务器类型
	Msgcenter *MsgCenter //消息中心
	si        *ServerIndex
}

func (ss *ServerService) Get(id uint32) *Server {
	return ss.si.get(id)
}

// Forward rpc.Forward调用实现,解析转发来的消息
func (ss *ServerService) Forward(stream Server_ForwardServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("invalid rpc,no metadata")
	}
	clientid, ok := md["clientid"]
	if !ok {
		return fmt.Errorf("invalid client,no clientid")
	}
	//获得clientid
	cid, err := strconv.Atoi(clientid[0])
	if err != nil {
		return err
	}
	server := NewServer(ss.Typeid, stream, ss.Msgcenter, uint32(cid))
	err = ss.si.put(server)
	if err != nil {
		return err
	}
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
	ss.si = NewServerIndex()
	if err == nil {
		rpcs := grpc.NewServer()
		RegisterServerServer(rpcs, ss)
		return rpcs.Serve(lis)
	}
	return err
}
