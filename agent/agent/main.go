package main

import (
	"lemna/agent"
	"lemna/logger"
	"time"
)

var cm clientManager
var rb balancer
var t token

func main() {
	rb.start()
	cm.clients = make(map[int32]agent.Client)
	t.data = tokenMap
	as := agent.AgentService{Port: ":9999", Balancer: &rb, Cm: &cm, Token: &t}
	if err := as.Start(); err == nil {
		logger.Info("agent is  Running")
		for {
			time.Sleep(time.Second)
		}
	} else {
		logger.Error(err)
	}
}
