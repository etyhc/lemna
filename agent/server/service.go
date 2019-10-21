// Package server 实现了代理服务器的服务器端rpc服务.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"lemna/agent"
	"lemna/arpc"
	"lemna/logger"
	"net"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Service 代理rpc服务，接受服务器连接
type Service struct {
	addr   string //代理地址
	ctp    agent.TargetPool
	schers map[uint32]scher
}

func (s *Service) add(ser Server) {
	id := ser.ID()
	ls, ok := s.schers[id]
	if !ok {
		switch ser.info().Sche {
		case SERVERSCHENIL:
			s.schers[id] = newNilScher()
		case SERVERSCHELOAD:
			s.schers[id] = newLoadScher()
		case SERVERSCHEROUND:
			s.schers[id] = newRoundScher()
		}
		ls, ok = s.schers[id]
	}
	ls.add(ser)
}

func (s *Service) del(ser Server) {
	sche, ok := s.schers[ser.ID()]
	if ok {
		sche.del(ser)
	}
}

//GetTarget 服务器池接口实现
func (s *Service) GetTarget(target uint32) agent.Target {
	if scher, ok := s.schers[target]; ok {
		return scher.sche()
	}
	return nil
}

// NewService 新代理服务
func NewService(addr string) *Service {
	return &Service{addr: addr, schers: make(map[uint32]scher)}
}

//从header中获取server info
func getInfo(ctx context.Context) (Info, error) {
	var info Info
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return info, fmt.Errorf("invalid rpc,no metadata")
	}
	is, ok := md["info"]
	if !ok {
		return info, fmt.Errorf("invalid server, no info")
	}
	err := json.Unmarshal([]byte(is[0]), &info)
	return info, err
}

// Forward arpc.ArpcServer.Forward接口实现
func (s *Service) Forward(stream arpc.Srpc_ForwardServer) error {
	info, err := getInfo(stream.Context())
	if err != nil {
		return err
	}
	target := NewFTarget(stream, info)
	s.add(target)
	logger.Infof("Server%v Connected.", *target)
	return target.Forward(s.ctp)
}

// Multicast arpc.ArpcServer.Multicast接口实现
func (s *Service) Multicast(stream arpc.Srpc_MulticastServer) error {
	info, err := getInfo(stream.Context())
	if err != nil {
		return err
	}
	return NewMTarget(stream, info).Forward(s.ctp)
}

// Other arpc.ArpcServer.Other接口实现
func (s *Service) Other(stream arpc.Srpc_OtherServer) error {
	//TODO 待实现
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
	logger.Infof("Start server service at %s", s.addr)
	return gs.Serve(lis)
}
