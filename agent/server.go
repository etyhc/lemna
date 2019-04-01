package agent

import (
	"lemna/agent/rpc"
)

// Server 服务器流接口，agent会将用户消息通过Stream().Send()转发给此服务器
type Server interface {
	Stream() rpc.Server_ForwardClient
	TypeID() int32
}

// Balancer 服务负载均衡器
//          agent转发消息时,都从均衡器中得到一个可用的Server并将消息转发给此Server
type Balancer interface {
	GetServer(target int32, sessionid int32) (Server, bool)
}
