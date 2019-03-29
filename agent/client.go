package agent

import (
	"lemna/agent/rpc"
)

type Client interface {
	Stream() rpc.Client_ForwardServer
	SetStream(rpc.Client_ForwardServer)
	Run() error
	SessionID() int32
}

type ClientManager interface {
	GetClient(sessionid int32) (Client, error)
	DelClient(sessionid int32)
}
