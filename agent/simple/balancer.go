package simple

import (
	"lemna/agent"
	"lemna/rpc"
	"time"
)

type SimpleBalancer struct {
	servers map[int32]*rpc.ClientService
	as      *agent.AgentService
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
				sb.as.RunServer(cs)
				delete(sb.servers, cs.Typeid)
			}
			time.Sleep(time.Second)
		}
	}()
}

func (sb *SimpleBalancer) Start(as *agent.AgentService) {
	sb.servers = make(map[int32]*rpc.ClientService)
	sb.as = as
	sb.registerServer(&rpc.ClientService{Typeid: 1, Addr: ":10001"})
	sb.registerServer(&rpc.ClientService{Typeid: 2, Addr: ":10002"})
}
