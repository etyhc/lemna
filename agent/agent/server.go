package main

import (
	"io"
	"lemna/agent/rpc"
	"lemna/logger"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func (s *server) Stream() rpc.Server_ForwardClient {
	return s.stream
}

type server struct {
	Typeid int32
	Port   string
	stream rpc.Server_ForwardClient
}

func (s *server) init() error {
	conn, err := grpc.Dial(s.Port, grpc.WithInsecure())
	if err == nil {
		sc := rpc.NewServerClient(conn)
		s.stream, err = sc.Forward(context.Background())
		if err == nil || err == io.EOF {
			logger.Info(s.Typeid, " is ready.")
			return nil
		}
	}
	return err
}

func (s *server) run() {
	for {
		sfmsg, err := s.stream.Recv()
		if err != nil && err != io.EOF {
			break
		}
		if client, err := cm.GetClient(sfmsg.Target); err == nil {
			logger.Info("s to c ", sfmsg.Target)
			sfmsg.Target = s.Typeid
			client.Stream().Send(sfmsg)
		}
	}
}
