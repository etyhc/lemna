package simple

import (
	"lemna/agent"
	"time"
)

type SimpleBalancer struct {
	servers map[int32]*SimpleServer
	as      *agent.AgentService
}

func (sb *SimpleBalancer) GetServer(target int32, sessionid int32) (s agent.Server, ok bool) {
	s, ok = sb.servers[target]
	return
}

func (sb *SimpleBalancer) registerServer(s *SimpleServer) {
	go func() {
		for {
			if s.init() == nil {
				sb.servers[s.Typeid] = s
				sb.as.RunServer(s)
				delete(sb.servers, s.Typeid)
			}
			time.Sleep(time.Second)
		}
	}()
}

func (sb *SimpleBalancer) Start(as *agent.AgentService) {
	sb.servers = make(map[int32]*SimpleServer)
	sb.as = as
	sb.registerServer(&SimpleServer{Typeid: 1, Port: ":10001"})
	sb.registerServer(&SimpleServer{Typeid: 2, Port: ":10002"})
}
