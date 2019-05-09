// Package arpc 基于grpc的代理服务器的基础定义.
//
//             定义了2个grpc服务Outside和Inside
//             客户端<----------->代理<---------->服务器
//                    Outside rpc      Inside rpc
package arpc

import (
	fmt "fmt"
	"net"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// serverIndex 代理服务器索引
type serverIndex struct {
	index map[uint32]*Server
	mu    sync.Mutex
}

func (si *serverIndex) get(id uint32) *Server {
	si.mu.Lock()
	defer si.mu.Unlock()
	return si.index[id]
}

func (si *serverIndex) put(s *Server) error {
	si.mu.Lock()
	defer si.mu.Unlock()
	if _, ok := si.index[s.ID()]; ok {
		return fmt.Errorf("clientid<%d> conflcit", s.ID())
	}
	si.index[s.ID()] = s
	return nil
}

func (si *serverIndex) remove(s *Server) {
	si.mu.Lock()
	defer si.mu.Unlock()
	delete(si.index, s.ID())
}

// ServerService 服务器rpc服务，用于服务器开发
//               服务接收代理服务器连接，并执行代理服务器的rpc调用
type ServerService struct {
	addr   string       //服务器地址
	typeid uint32       //服务器类型
	mc     *MsgCenter   //消息中心
	si     *serverIndex //代理服务器索引
}

// NewServerService 新rpc服务
func NewServerService(addr string, typeid uint32, mc *MsgCenter) *ServerService {
	return &ServerService{
		addr:   addr,
		typeid: typeid,
		mc:     mc,
		si:     &serverIndex{index: make(map[uint32]*Server)}}
}

// Get 根据唯一ID得到代理服务器rpc服务端
func (ss *ServerService) Get(id uint32) *Server {
	return ss.si.get(id)
}

// Forward rpc.Forward调用实现
//         执行代理服务器的Forward调用
//         解析代理转发来的消息
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
	server := NewServer(ss.typeid, stream, ss.mc, uint32(cid))
	err = ss.si.put(server)
	if err != nil {
		return err
	}
	defer ss.si.remove(server)
	return server.Handle()
}

// Run 运行rpc服务,阻塞的
func (ss *ServerService) Run() error {
	lis, err := net.Listen("tcp", ss.addr)
	if err == nil {
		rpcs := grpc.NewServer()
		RegisterServerServer(rpcs, ss)
		return rpcs.Serve(lis)
	}
	return err
}
