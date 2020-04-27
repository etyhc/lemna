package srpc

import (
	context "context"
	"lemna/logger"
	"net"

	grpc "google.golang.org/grpc"
)

//Service 服务器调用服务
type Service struct {
	addr string
}

// NewService 服务器调用服务
func NewService(addr string) *Service {
	return &Service{addr: addr}
}

//Run 运行服务
func (s *Service) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	gs := grpc.NewServer()
	RegisterStoSServer(gs, s)
	logger.Infof("Start server service at %s", s.addr)
	return gs.Serve(lis)
}

//Call StoSServer.Call 实现
func (s *Service) Call(context.Context, *CallMsg) (*CallMsg, error) {
	//TODO 处理来自服务器的调用
	return nil, nil
}
