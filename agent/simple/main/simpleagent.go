package main

import (
	"lemna/agent"
	"lemna/agent/simple"
	"lemna/logger"
)

var cm *simple.SimpleClientManager
var rb *simple.SimpleBalancer
var t *simple.SimpleToken
var as *agent.Service

func init() {
	cm = simple.NewSimpleClientManager()
	t = simple.NewSimpleToken()
	rb = simple.NewSimpleBalancer()
	as = &agent.Service{Port: ":9999", Balancer: rb, Cm: cm, Token: t}
	rb.Start(as)
}
func main() {
	logger.SetLevel(logger.DEBUG)
	if err := as.Run(); err != nil {
		logger.Error(err)
	}
}
