package agent

import (
	"lemna/rpc"
)

// Client 用户流接口,agent会将服务器消息通过Stream().Send()转发给此用户
type Client interface {
	Stream() rpc.Client_ForwardServer
	SetStream(rpc.Client_ForwardServer)
	SessionID() int32
	Error(interface{}) error
}

// ClientManager 用户管理器接口
//               管理器应负责此网关所有用户的注册，生成，管理，注销,资源回收
type ClientManager interface {
	NewClient(sessionid int32) (Client, error)
	GetClient(sessionid int32) (Client, error)
	DelClient(sessionid int32)
}
