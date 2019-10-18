package client

import (
	fmt "fmt"
	"lemna/agent"
	"lemna/arpc"
	"lemna/logger"
)

// Target 客户端rpc调用服务
type Target struct {
	stream arpc.Crpc_ForwardServer //网络流
	uid    uint32                  //客户端id
	cache  map[uint32]agent.Target //转发目标缓存
}

// NewTarget 新客户端
//         s 客户端网络流
//        id 客户端uid，客户端唯一
func newTarget(s arpc.Crpc_ForwardServer, id uint32) *Target {
	return &Target{stream: s, uid: id, cache: make(map[uint32]agent.Target)}
}

// ID 客户端UID
func (t *Target) ID() uint32 {
	return t.uid
}

func (t *Target) wraperr(err error) error {
	return fmt.Errorf("<uid=%d>%w", t.uid, err)
}

// Send 向客户端发送转发消息
func (t *Target) Send(msg *arpc.ForwardMsg) error {
	return t.stream.Send(msg)
}

// Forward 接收客户端转发消息
func (t *Target) Forward(pool agent.TargetPool) error {
	for {
		fmsg, err := t.stream.Recv()
		if err != nil {
			for _, server := range t.cache {
				agent.InvalidTarget(server, t.uid)
			}
			return t.wraperr(err)
		}

		server, isCached := t.cache[fmsg.Target]
		if !isCached {
			server := pool.GetTarget(fmsg.Target)
			if server == nil { //目标服务器无效,丢弃这次数据
				agent.InvalidTarget(t, fmsg.Target)
				logger.Errorf("not find server<%d>", fmsg.Target)
				continue
			}
		}

		//转发指令
		err = agent.T2T(t, server, fmsg)
		if err != nil { //转发失败
			logger.Error(err)
			if isCached {
				delete(t.cache, server.ID())
			}
		} else {
			if !isCached {
				t.cache[server.ID()] = server
			}
		}
	}
}

//Bind 绑定服务器，默认此服务器提供服务
func (t *Target) Bind(server agent.Target) {
	t.cache[server.ID()] = server
}
