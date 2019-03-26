package main

import (
	"lemna/agent"
	"time"
)

func main() {
	as := agent.AgentService{Port: ":9999"}
	as.RegisterServer(&agent.Server{Typeid: 1, Port: ":10001"})
	as.RegisterServer(&agent.Server{Typeid: 2, Port: ":10000"})
	as.Start()
	for {
		time.Sleep(time.Second)
	}
}
