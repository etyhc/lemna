package main

import (
	"flag"
	"lemna/agent"
	"lemna/agent/impl"
	configrpc "lemna/config/rpc"
	"lemna/logger"
)

var as *agent.Service

func init() {
	addr = flag.String("addr", ":9999", "要绑定的地址")
	configaddr = flag.String("config", configrpc.ConfigServerAddr, "订阅服务器地址")
	h = flag.Bool("h", false, "this help")
	bp = impl.NewBalancePool()
}

var bp *impl.BalancePool
var addr *string
var configaddr *string
var h *bool

func main() {
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	if err := bp.Start(*configaddr); err != nil {
		logger.Error(err)
		return
	}
	as = agent.NewService(*addr, bp, impl.NewSimpleToken())
	if err := as.Run(); err != nil {
		logger.Error(err)
	}
}
