package simple

import (
	"io"
	"lemna/agent/rpc"
	"lemna/logger"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func (s *SimpleServer) Stream() rpc.Server_ForwardClient {
	return s.stream
}

func (s *SimpleServer) TypeID() int32 {
	return s.Typeid
}

type SimpleServer struct {
	Typeid int32
	Port   string
	stream rpc.Server_ForwardClient
}

func (s *SimpleServer) init() error {
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
