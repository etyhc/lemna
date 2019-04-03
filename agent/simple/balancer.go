package simple

import (
	"lemna/agent"
	"lemna/logger"
	"lemna/rpc"
	"time"
)

type SimpleBalancer struct {
	servers map[int32]*rpc.ClientService
	as      *agent.Service
}

func (sb *SimpleBalancer) GetServer(target int32, sessionid int32) (s agent.Server, ok bool) {
	s, ok = sb.servers[target]
	return
}

func (sb *SimpleBalancer) registerServer(cs *rpc.ClientService) {
	go func() {
		for {
			if cs.Init() == nil {
				sb.servers[cs.Typeid] = cs
				logger.Infof("%d is running", cs.Typeid)
				err := sb.as.RunServer(cs)
				logger.Error(err)
				delete(sb.servers, cs.Typeid)
			}
			time.Sleep(time.Second)
		}
	}()
}

func NewSimpleBalancer() (sb *SimpleBalancer) {
	sb = &SimpleBalancer{}
	sb.servers = make(map[int32]*rpc.ClientService)
	return
}

func (sb *SimpleBalancer) Start(as *agent.Service) {
	sb.as = as
	sb.registerServer(&rpc.ClientService{Typeid: 1, Addr: ":10001"})
	sb.registerServer(&rpc.ClientService{Typeid: 2, Addr: ":10002"})
}
