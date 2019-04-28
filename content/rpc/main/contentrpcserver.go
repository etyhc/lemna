package main

import (
	"flag"
	"lemna/content/rpc"
	"lemna/logger"
)

var addr *string
var h *bool

func init() {
	addr = flag.String("addr", rpc.SERVERADDR, "要绑定的地址")
	h = flag.Bool("h", false, "this help")
}

func main() {
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	rcs := rpc.NewChannelService(*addr)
	err := rcs.Run()
	if err != nil {
		logger.Error(err)
	}
}
