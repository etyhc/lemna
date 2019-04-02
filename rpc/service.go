package rpc

import (
	"lemna/logger"
	"net"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ServerService rpc服务器端封装
type ServerService struct {
	Addr      string     //服务器地址
	Typeid    int32      //服务器类型
	Msgcenter *MsgCenter //消息中心
}

// Forward rpc.Forward调用实现,解析转发来的消息
func (ss *ServerService) Forward(stream Server_ForwardServer) (err error) {
	var in *ForwardMsg
	for {
		in, err = stream.Recv()
		if err == nil {
			err = ss.Msgcenter.Handle(in, stream)
		}
		if err != nil {
			return
		}
	}
}

// Run 运行rpc服务,阻塞的
func (ss *ServerService) Run() error {
	lis, err := net.Listen("tcp", ss.Addr)
	if err == nil {
		rpcs := grpc.NewServer()
		RegisterServerServer(rpcs, ss)
		return rpcs.Serve(lis)
	}
	return err
}

// ClientService rpc客户端服务
type ClientService struct {
	Addr   string //服务器地址
	Typeid int32  //服务器类型
	stream Server_ForwardClient
}

// Stream 服务器流
func (cs *ClientService) Stream() Server_ForwardClient {
	return cs.stream
}

// TypeID 服务器类型
func (cs *ClientService) TypeID() int32 {
	return cs.Typeid
}

// Init 初始化
func (cs *ClientService) Init() error {
	conn, err := grpc.Dial(cs.Addr, grpc.WithInsecure())
	if err == nil {
		sc := NewServerClient(conn)
		cs.stream, err = sc.Forward(context.Background())
		if err == nil {
			logger.Info(cs.Typeid, " is ready.")
		}
	}
	return err
}
