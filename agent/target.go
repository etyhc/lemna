package agent

import (
	"lemna/arpc"
)

// Target 转发目标，用于具体消息转发工作
type Target interface {
	// Send 向目标发送转发消息
	Send(*arpc.ForwardMsg) error
	Forward(TargetPool) error
	// ID 目标ID，客户端是唯一的，服务器端只是代表服务器类型
	//    代理连接很多同类型服务器，但同类型服务器只有一个为客户端服务
	ID() uint32
	Error(error) error
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
		return dest.Error(err)
	}
	return nil
}

// C2S 阻塞循环接收客户端消息并转发给相应服务器
//     在自己的缓存未找到转发服务器，再从转发服务器池寻找转发服务器
//     转发失败清除自己缓存的转发服务器，回复InvalidTargetMsg消息给客户端
//   src 客户端，消息来源
//  pool 服务器池，缓冲无服务器会从服务器池中获取服务器
//func C2S(src CTarget, pool TargetPool) error {
//	for {
//		fmsg, err := src.Recv()
//		if err != nil { //接收客户端数据失败，通知已缓冲服务器，客户端接收数据失败
//			for _, server := range src.Cache() {
//				invalidTarget(server, src.ID())
//			}
//			return src.Error(err) //结束客户端转发服务
//		}
//
//		dest, ok := src.Cache()[fmsg.Target]
//		if !ok {
//			d := pool.GetTarget(fmsg.Target)
//			if d == nil { //目标服务器无效,丢弃这次数据
//				invalidTarget(src, fmsg.Target)
//				logger.Errorf("not find server<%d>", fmsg.Target)
//				continue
//			}
//			dest = d.(STarget)
//		}
//
//		//转发指令
//		err = t2t(src, dest, fmsg)
//		if err != nil { //转发失败
//			logger.Error(err)
//			delete(src.Cache(), dest.ID())
//			continue
//		}
//	}
//}
//
//func S2C(src STarget, pool TargetPool) error {
//	for {
//		fmsg, err := src.Recv()
//		if err != nil { //结束服务器转发
//			return src.Error(err)
//		}
//
//		dest := pool.GetTarget(fmsg.Target)
//		if dest == nil {
//			invalidTarget(src, fmsg.Target)
//			logger.Errorf("not find client<%d>", fmsg.Target)
//			continue
//		}
//
//		//转发指令
//		err = t2t(src, dest, fmsg)
//		if err != nil { //转发失败
//			invalidTarget(src, fmsg.Target)
//			logger.Error(err)
//		}
//	}
//}
//
//// M2C 阻塞循环接收服务器消息并转发给客户端.
////   src 服务器目标，用于接收消息
////  pool 客户端池
//func M2C(src STarget, pool TargetPool) error {
//	for {
//		bmsg, err := src.Recvm()
//		if err != nil {
//			return src.Error(err)
//		}
//
//		for _, cid := range bmsg.Targets {
//			c := pool.GetTarget(cid)
//			if c == nil {
//				invalidTarget(src, cid)
//				logger.Errorf("not find client<%d>", cid)
//				continue
//			}
//
//			//转发指令
//			err = t2t(src, c, &arpc.ForwardMsg{Target: src.ID(), Msg: bmsg.Msg})
//			//转发失败
//			if err != nil {
//				logger.Error(err)
//			} else { //转发成功，Cache服务器
//				c.(CTarget).Cache()[src.ID()] = src
//			}
//		}
//	}
//}
