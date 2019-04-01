package main

import (
	"lemna/agent"
	"lemna/agent/simple"
	"lemna/logger"
	"time"
)

var cm simple.SimpleClientManager
var rb simple.SimpleBalancer
var t simple.SimpleToken
var as agent.AgentService

func main() {
	cm.Init()
	t.Init()
	as = agent.AgentService{Port: ":9999", Balancer: &rb, Cm: &cm, Token: &t}
	rb.Start(&as)
	if err := as.Start(); err == nil {
		logger.Info("agent is  Running")
		for {
			time.Sleep(time.Second)
		}
	} else {
		logger.Error(err)
	}
}
