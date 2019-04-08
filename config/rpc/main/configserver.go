package main

import (
	"flag"
	"lemna/config/rpc"
)

var addr *string
var h *bool

func init() {
	addr = flag.String("addr", rpc.ConfigServerAddr, "要绑定的地址")
	h = flag.Bool("h", false, "this help")
}

func main() {
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	rcs := rpc.NewChannelService(*addr)
	rcs.Run()
}
