package main

import (
	"lemna/agent"
	"lemna/agent/simple"
	"lemna/logger"
)

var as *agent.Service

func init() {
	as = &agent.Service{
		Port:     ":9999",
		Balancer: simple.NewSimpleBalancer(),
		Cm:       simple.NewSimpleClientManager(),
		Token:    simple.NewSimpleToken()}
}
func main() {
	if err := as.Run(); err != nil {
		logger.Error(err)
	}
}
