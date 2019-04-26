package main

import (
	"flag"
	"lemna/agent"
	"lemna/agent/impl"
	configrpc "lemna/config/rpc"
	"lemna/logger"
)

var cs *agent.ClientService

func init() {
	addr = flag.String("addr", ":9999", "要绑定的地址")
	configaddr = flag.String("config", configrpc.ConfigServerAddr, "订阅服务器地址")
	h = flag.Bool("h", false, "this help")
}

var sp impl.SubServerPool
var addr *string
var configaddr *string
var h *bool

func main() {
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	if err := sp.SubscribeServer(*configaddr); err != nil {
		logger.Error(err)
		return
	}
	cs = agent.NewClientService(*addr, &sp, impl.NewSimpleToken())
	if err := cs.Run(); err != nil {
		logger.Error(err)
	}
}
