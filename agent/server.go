package agent

import (
	"io"
	"log"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type Server struct {
	Typeid int32
	Port   string
	stream Server_ForwardClient
}

func (s *Server) init() error {
	conn, err := grpc.Dial(s.Port, grpc.WithInsecure())
	if err == nil {
		sc := NewServerClient(conn)
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
		msg, err := s.stream.Recv()
		if err != nil && err != io.EOF {
			break
		}
		if client, ok := clientMap[msg.ClientID]; ok {
			log.Println("s to c ", msg.ClientID)
			client.stream.Send(&ClientFwdMsg{ServerType: s.Typeid, Msg: msg.Msg})
		}
	}
}
