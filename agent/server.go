package agent

import (
	"io"
	"lemna/agent/rpc"
	"log"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type Server struct {
	Typeid int32
	Port   string
	stream rpc.Server_ForwardClient
}

func (s *Server) init() error {
	conn, err := grpc.Dial(s.Port, grpc.WithInsecure())
	if err == nil {
		sc := rpc.NewServerClient(conn)
		s.stream, err = sc.Forward(context.Background())
		if err == nil || err == io.EOF {
			log.Println(s.Typeid, " is ready.")
			return nil
		}
	}
	return err
}

func (s *Server) run() {
	for {
		sfmsg, err := s.stream.Recv()
		if err != nil && err != io.EOF {
			break
		}
		if client, ok := clientMap[sfmsg.Target]; ok {
			log.Println("s to c ", sfmsg.Target)
			sfmsg.Target = s.Typeid
			client.stream.Send(sfmsg)
		}
	}
}
