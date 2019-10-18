// Package server 实现了代理服务器的服务器端rpc服务.
package server

import (
	"lemna/agent"
	"lemna/arpc"
	"net"

	grpc "google.golang.org/grpc"
)

// Service 代理rpc服务，接受服务器连接
type Service struct {
	addr   string //代理地址
	ctp    agent.TargetPool
	schers map[uint32]scher
}

type scher interface {
	get(uint32) *Server
	sche() *Server
	del(uint32)
	up(*Server)
}

// GetTarget 服务器池接口实现
func (s *Service) GetTarget(target uint32) agent.Target {
	if scher, ok := s.schers[target]; ok {
		server := scher.sche()
		if server != nil {
			return server.target
		}
	}
	return nil
}

// NewService 新代理服务
func NewService(addr string) *Service {
	return &Service{addr: addr, schers: make(map[uint32]scher)}
}

func (s *Service) Update(server *Server) {
	scher, isInited := s.schers[server.target.ID()]
	if !isInited {
		//TODO 初始化一个调度器
	}
	scher.up(server)
}

// Forward arpc.ArpcServer.Forward接口实现
func (s *Service) Forward(stream arpc.Srpc_ForwardServer) error {
	var info Info
	ft := NewFTarget(stream, &info)
	return ft.Forward(s.ctp)
}

// Multicast arpc.ArpcServer.Multicast接口实现
func (s *Service) Multicast(stream arpc.Srpc_MulticastServer) error {
	return NewMTarget(stream).Forward(s.ctp)
}

// Multicast arpc.ArpcServer.Multicast接口实现
func (s *Service) Other(stream arpc.Srpc_OtherServer) error {
	return nil
}

// Run 运行代理服务,接受服务器的连接
func (s *Service) Run(ctp agent.TargetPool) error {
	s.ctp = ctp
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	gs := grpc.NewServer()
	arpc.RegisterSrpcServer(gs, s)
	return gs.Serve(lis)
}
