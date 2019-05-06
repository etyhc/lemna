package main

import (
	"flag"
	"lemna/agent"
	"lemna/agent/client"
	"lemna/agent/server"
	"lemna/content/crpc"
	"lemna/logger"
)

func init() {
	addr = flag.String("addr", ":9999", "要绑定的地址")
	configaddr = flag.String("config", crpc.SERVERADDR, "订阅服务器地址")
	h = flag.Bool("h", false, "this help")
}

var addr *string
var configaddr *string
var h *bool

func main() {
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	sp := server.NewSubServerPool(*configaddr)
	cp := client.NewClientService(*addr, client.NewSimpleToken())
	as := agent.NewService(sp, cp)
	if err := as.Run(); err != nil {
		logger.Error(err)
	}
}
