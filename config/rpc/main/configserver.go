package main

import (
	"lemna/config/rpc"
)

func main() {
	rcs := rpc.NewService(rpc.ConfigServerAddr)
	rcs.Run()
}
