package impl

import (
	"lemna/agent"
)

// NilScheduler 代理服务器不提供调度服务,所以返回nil
//              应该服务器先向客户端发送消息，代理自动缓存服务器，建立服务器和客户端的对应关系
type NilScheduler struct {
}

func (bs *NilScheduler) schedule(servers map[string]*agent.Server, client *agent.Client) *agent.Server {
	return nil
}
