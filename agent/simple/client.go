package simple

import (
	fmt "fmt"
	"lemna/agent/rpc"
)

type SimpleClient struct {
	sessionid int32
	stream    rpc.Client_ForwardServer
}

func (c *SimpleClient) Stream() rpc.Client_ForwardServer {
	return c.stream
}

func (c *SimpleClient) SessionID() int32 {
	return c.sessionid
}

func (c *SimpleClient) session() string {
	return fmt.Sprintf("%d", c.sessionid)
}

func (c *SimpleClient) SetStream(stream rpc.Client_ForwardServer) {
	c.stream = stream
}
