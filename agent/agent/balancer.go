package main

import (
	"lemna/agent"
	"time"
)

type balancer struct {
	servers map[int32]*server
}

func (rb *balancer) GetServer(target int32, clientid int32) (s agent.Server, ok bool) {
	s, ok = rb.servers[target]
	return
}

func (rb *balancer) registerServer(s *server) {
	go func() {
		for {
			if s.init() == nil {
				rb.servers[s.Typeid] = s
				s.run()
				delete(rb.servers, s.Typeid)
			}
			time.Sleep(time.Second)
		}
	}()
}

func (rb *balancer) start() {
	rb.servers = make(map[int32]*server)
	rb.registerServer(&server{Typeid: 1, Port: ":10001"})
	rb.registerServer(&server{Typeid: 2, Port: ":10002"})
}
