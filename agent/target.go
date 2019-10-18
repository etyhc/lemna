package agent

import (
	"lemna/arpc"
)

// Target 转发目标，用于具体消息转发工作
type Target interface {
	//Send 向目标发送转发消息
	Send(*arpc.ForwardMsg) error
	// ID 目标ID，客户端是UID，服务器端只是代表服务器类型
	//    代理连接很多同类型服务器，但同类型服务器只有一个为客户端服务
	ID() uint32
}

// TargetPool 目标池
type TargetPool interface {
	// GetTarget 得到目标，无目标返回nil
	GetTarget(uint32) Target
	Run(TargetPool) error
}

// InvalidTarget 通知src，dest是无效目标
func InvalidTarget(src Target, dest uint32) {
	_ = src.Send(MsgWrapper(msgkits.Builder.InvalidTarget(), dest))
}

// T2T 转发功能封装
func T2T(src Target, dest Target, msg *arpc.ForwardMsg) error {
	msg.Target = src.ID()
	err := dest.Send(msg)
	if err != nil {
		InvalidTarget(src, dest.ID())
	}
	return err
}
