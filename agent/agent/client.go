package main

import (
	fmt "fmt"
	"lemna/agent/rpc"
)

type client struct {
	sessionid int32
	stream    rpc.Client_ForwardServer
}

func (c *client) Stream() rpc.Client_ForwardServer {
	return c.stream
}

func (c *client) SessionID() int32 {
	return c.sessionid
}

func (c *client) session() string {
	return fmt.Sprintf("%d", c.sessionid)
}

func (c *client) SetStream(stream rpc.Client_ForwardServer) {
	c.stream = stream
}

func (c *client) Run() error {
	for {
		cfmsg, err := c.stream.Recv()
		if err == nil {
			//转发指令
			if server, ok := rb.GetServer(cfmsg.Target, c.sessionid); ok {
				cfmsg.Target = c.sessionid
				err = server.Stream().Send(cfmsg)
			} else {
				err = fmt.Errorf("invaild servertype %d", cfmsg.Target)
			}
		}
		if err != nil {
			return err
		}
	}
}
