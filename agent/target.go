package agent

import (
	"lemna/agent/arpc"
	"lemna/logger"
)

// Target 转发目标，用于具体消息转发工作
type Target interface {
	// Send 向目标发送转发消息
	Send(*arpc.ForwardMsg) error
	// ID 目标ID，客户端是唯一的，服务器端只是代表服务器类型
	//    代理连接很多同类型服务器，但同类型服务器只有一个为客户端服务
	ID() int32
	// Error 为err添加目标信息
	Error(err interface{}) error
}

// CTarget 客户端目标
//         客户端只有转发消息
//         客户端有服务器端缓存,最近一次与客户端交互服务器会被缓存
//         客户端消息会转发给缓存中的服务器，如果缓存没有才会在服务器池中寻找服务器转发
type CTarget interface {
	Target
	Recv() (*arpc.ForwardMsg, error)
	//服务器缓存
	Cache() map[int32]STarget
}

// STarget 服务器目标
//         服务器只有广播消息，代理接收到广播消息会将广播消息重新封装成转发消息，并转发给所有要广播的客户端
type STarget interface {
	Target
	Recv() (*arpc.BroadcastMsg, error)
}

// TargetPool 目标池
type TargetPool interface {
	// GetTarget 得到目标，无目标返回nil
	GetTarget(int32) Target
	// 绑定转发池
	Bind(TargetPool)
	// 运行，阻塞的
	Run() error
}

// invalidTarget 将无效目标dest，发给src
func invalidTarget(src Target, dest int32) {
	itm, err := arpc.WrapFMNoCheck(dest, &InvalidTargetMsg{})
	if err == nil {
		src.Send(itm)
	}
}

// C2S 阻塞循环接收客户端消息并转发给相应服务器
//   src 客户端，消息来源
//  pool 服务器池，缓冲无服务器会从服务器池中获取服务器
func C2S(src CTarget, pool TargetPool) error {
	for {
		fmsg, err := src.Recv()
		if err != nil {
			//Say Bye
			for _, server := range src.Cache() {
				invalidTarget(server, src.ID())
			}
			return src.Error(err)
		}

		dest, ok := src.Cache()[fmsg.Target]
		if !ok {
			dest = pool.GetTarget(fmsg.Target).(STarget)
			if dest == nil {
				invalidTarget(src, fmsg.Target)
				logger.Errorf("not find server<%d>", fmsg.Target)
				continue
			}
		}

		//转发指令
		fmsg.Target = src.ID()
		err = dest.Send(fmsg)
		//转发失败
		if err != nil {
			invalidTarget(src, fmsg.Target)
			logger.Error(dest.Error(err))
			delete(src.Cache(), dest.ID())
			continue
		}
	}
}

// S2C 阻塞循环接收服务器消息并转发给客户端.
//   src 服务器目标，用于接收消息
//  pool 客户端池
func S2C(src STarget, pool TargetPool) error {
	for {
		bmsg, err := src.Recv()
		if err != nil {
			return src.Error(err)
		}

		for _, cid := range bmsg.Targets {
			c := pool.GetTarget(cid).(CTarget)
			if c == nil {
				invalidTarget(src, cid)
				logger.Errorf("not find client<%d>", cid)
				continue
			}

			//转发指令
			err = c.Send(&arpc.ForwardMsg{Target: src.ID(), Msg: bmsg.Msg})
			//转发失败
			if err != nil {
				invalidTarget(src, cid)
				logger.Error(c.Error(err))
				delete(c.Cache(), src.ID())
			} else {
				c.Cache()[src.ID()] = src
			}
		}
	}
}
