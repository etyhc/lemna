package agent

import "lemna/agent/rpc"

// Server 服务器流接口，agent会将用户消息通过Stream().Send()转发给此服务器
type Server interface {
	Stream() rpc.Server_ForwardClient
	TypeID() int32
	Error(interface{}) error
}

// Balancer 服务负载均衡器
//          代理转发消息时,从均衡器中得到一个可用的Server并将消息转发给此Server
//          均衡器应负责Server的发现，管理，负载均衡，注销，资源回收
//          均衡器应在Server就绪时调用AgentService.RunServer让Server开始转发工作
type Balancer interface {
	GetServer(target int32, client Client) (Server, bool)
}
