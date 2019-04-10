package main

import (
	"lemna/agent"
	"lemna/agent/impl"
	"lemna/logger"
)

var as *agent.Service

func init() {
	bp := impl.NewBalancePool()
	bp.Start()
	as = agent.NewService(":9999", bp, impl.NewSimpleToken())
}

func main() {
	if err := as.Run(); err != nil {
		logger.Error(err)
	}
}
