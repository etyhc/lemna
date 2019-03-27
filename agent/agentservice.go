package agent

import (
	"errors"
	fmt "fmt"
	"lemna/agent/rpc"
	"log"
	"net"
	"strconv"
	"time"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AgentService ...
type AgentService struct {
	Port string
}

var clientMap = make(map[int32]*client)
var serverMap = make(map[int32]*Server)
var tokenMap = map[string]int32{"token": 1}

func (as *AgentService) Register(cont context.Context, msg *rpc.ClientRegMsg) (*rpc.ClientRegMsg, error) {
	if sessionid, ok := tokenMap[msg.Token]; ok {
		if _, ok := clientMap[sessionid]; !ok {
			clientMap[sessionid] = &client{id: sessionid, login: true}
		}
		grpc.SetHeader(cont, metadata.Pairs("sessionid", clientMap[sessionid].session()))
		return msg, nil
	}
	return nil, fmt.Errorf("invalid token %s", msg.Token)
}

func (as *AgentService) Forward(stream rpc.Client_ForwardServer) error {

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		if session, ok := md["sessionid"]; ok {
			tmp, _ := strconv.Atoi(session[0])
			sessionid := int32(tmp)
			if client, ok := clientMap[sessionid]; ok {
				client.stream = stream
				err := client.run()
				delete(clientMap, sessionid)
				return err
			}
		}
	}
	return errors.New("invalid client")
}

func (as *AgentService) Start() {
	lis, err := net.Listen("tcp", as.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rpc.RegisterClientServer(s, as)
	log.Println("agent running")
	s.Serve(lis)
}

func (as *AgentService) RegisterServer(s *Server) {
	go func() {
		for {
			if s.init() == nil {
				serverMap[s.Typeid] = s
				s.run()
				delete(serverMap, s.Typeid)
			}
			time.Sleep(time.Second)
		}
	}()
}
