package main

import (
	"flag"
	"lemna/agent"
	"lemna/agent/client"
	"lemna/agent/server"
	"lemna/logger"
)

func init() {
	caddr = flag.String("caddr", ":9999", "crpc地址要绑定的地址")
	saddr = flag.String("saddr", ":10000", "srpc地址要绑定的地址")
	h = flag.Bool("h", false, "this help")
}

var caddr *string
var saddr *string
var h *bool

func main() {
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	sp := server.NewService(*saddr)
	cp := client.NewService(*caddr, client.NewSimpleToken())
	as := agent.NewService(sp, cp)
	if err := as.Run(); err != nil {
		logger.Error(err)
	}
}
