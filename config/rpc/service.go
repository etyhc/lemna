package rpc

import (
	"lemna/logger"
	"net"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

var ConfigServerAddr = ":10000"

type Service struct {
	subscribers []Config_SubscribeServer
	servers     map[string]int
	addr        string
}

func NewService(addr string) *Service {
	return &Service{
		subscribers: []Config_SubscribeServer{},
		servers:     make(map[string]int),
		addr:        addr}
}

func (s *Service) Publish(ctx context.Context, msg *ConfigMsg) (*ConfigMsg, error) {
	s.servers[msg.Info] = s.servers[msg.Info] + 1
	logger.Debug("pub: ", msg.Info)
	alive := s.subscribers[:0]
	for _, stream := range s.subscribers {
		err := stream.Send(msg)
		if err == nil {
			logger.Debug("sub: ...")
			alive = append(alive, stream)
		} else {
			logger.Error(err)
		}
	}
	s.subscribers = alive
	return msg, nil
}

func (s *Service) Subscribe(msg *ConfigMsg, stream Config_SubscribeServer) error {
	for info, _ := range s.servers {
		msg.Info = info
		err := stream.Send(msg)
		if err != nil {
			return err
		}
		logger.Debug("sub: ", msg.Info)
	}
	s.subscribers = append(s.subscribers, stream)
	select {
	case <-stream.Context().Done():
	}
	return nil
}

func (fs *Service) Run() error {
	lis, err := net.Listen("tcp", fs.addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	RegisterConfigServer(s, fs)
	return s.Serve(lis)
}
