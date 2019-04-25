package impl

import (
	"lemna/agent"
)

type BindScheduler struct {
}

// schedule 客户端不应该主动绑定服务器,所以返回nil
//          应该服务器先向客户端发送消息，客户端自动缓存自己，知道是哪个服务器为他服务
func (bs *BindScheduler) schedule(servers map[string]*server, client *agent.Target) *agent.Target {
	return nil
}
