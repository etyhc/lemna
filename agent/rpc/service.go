package rpc

import (
	fmt "fmt"
	"net"

	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Service struct {
	Addr      string     //服务器地址
	Typeid    int32      //服务器类型
	Msgcenter *MsgCenter //消息中心
}

// ServerService rpc服务器端封装，用于服务器开发
type ServerService struct {
	Addr      string     //服务器地址
	Typeid    int32      //服务器类型
	Msgcenter *MsgCenter //消息中心
	stream    Server_ForwardServer
}

// Forward rpc.Forward调用实现,解析转发来的消息
func (ss *ServerService) Forward(stream Server_ForwardServer) error {
	ss.stream = stream
	for {
		in, err := stream.Recv()
		if err == nil {
			err = ss.Msgcenter.Handle(in)
		}
		if err != nil {
			return err
		}
	}
}

func (ss *ServerService) Broadcast(targets []int32, msg proto.Message) error {
	sendmsg, err := ss.Msgcenter.WrapBroadcast(targets, msg)
	if err == nil {
		return ss.stream.Send(sendmsg)
	}
	return err
}

func (ss *ServerService) Send(target int32, msg proto.Message) error {
	sendmsg, err := ss.Msgcenter.WrapBroadcast([]int32{target}, msg)
	if err == nil {
		return ss.stream.Send(sendmsg)
	}
	return err
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

// ClientService rpc客户端服务封装，用于客户端开发
type ClientService struct {
	Addr      string     //服务器地址
	Typeid    int32      //服务器类型
	Msgcenter *MsgCenter //消息中心
	forwarder Server_ForwardClient
}

// Forwarder 服务器转发器
func (cs *ClientService) Forwarder() Server_ForwardClient {
	return cs.forwarder
}

// TypeID 服务器类型
func (cs *ClientService) TypeID() int32 {
	return cs.Typeid
}

func (cs *ClientService) Error(err interface{}) error {
	return fmt.Errorf("<%s> %s", cs.Addr, err)
}

// Init 初始化
func (cs *ClientService) Init() error {
	conn, err := grpc.Dial(cs.Addr, grpc.WithInsecure())
	if err == nil {
		sc := NewServerClient(conn)
		cs.forwarder, err = sc.Forward(context.Background())
	}
	return err
}
