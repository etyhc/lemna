package main

import (
	"flag"
	"lemna/content/crpc"
	"lemna/logger"
)

var addr *string
var h *bool

func init() {
	addr = flag.String("addr", crpc.SERVERADDR, "要绑定的地址")
	h = flag.Bool("h", false, "this help")
}

func main() {
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	rcs := crpc.NewChannelService(*addr)
	err := rcs.Run()
	if err != nil {
		logger.Error(err)
	}
}
