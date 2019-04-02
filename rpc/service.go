package rpc

import (
	"lemna/logger"
	"net"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ServerService struct {
	Addr   string
	Typeid int32
}

func (ss *ServerService) Forward(stream Server_ForwardServer) (err error) {
	var in *ForwardMsg
	for {
		in, err = stream.Recv()
		if err == nil {
			err = ForwardMsgHandle(in, stream)
		}
		if err != nil {
			return
		}
	}
}

func (ss *ServerService) Run() error {
	lis, err := net.Listen("tcp", ss.Addr)
	if err == nil {
		rpcs := grpc.NewServer()
		RegisterServerServer(rpcs, ss)
		return rpcs.Serve(lis)
	}
	return err
}

type ClientService struct {
	Addr   string
	Typeid int32
	stream Server_ForwardClient
}

func (cs *ClientService) Stream() Server_ForwardClient {
	return cs.stream
}

func (cs *ClientService) TypeID() int32 {
	return cs.Typeid
}

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
