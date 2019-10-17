package server

import (
	"fmt"
	"lemna/agent"
	"lemna/arpc"
	"lemna/logger"
)

// FTarget rpc代理端(代理服务器作为rpc服务器)
type FTarget struct {
	afs   arpc.Arpc_ForwardServer //Forward调用接收、发送端
	info  *Info                   //服务器信息
	Round uint32                  //服务器被调度次数
}

// NewFTarget 新服务器
//    client rpc客户端
//      info 订阅的服务器信息
func NewFTarget(afs arpc.Arpc_ForwardServer, info *Info) *FTarget {
	return &FTarget{afs: afs, info: info}
}

// Error 附加服务器信息到错误信息上
func (ft *FTarget) Error(err error) error {
	return fmt.Errorf("<type=%d>%w", ft.info.Type, err)
}

// Send 发送转发消息给服务器
func (ft *FTarget) Send(msg *arpc.ForwardMsg) error {
	return ft.afs.Send(msg)
}

// ID 服务器类型ID
func (ft *FTarget) ID() uint32 {
	return ft.info.Type
}

func (ft *FTarget) Forward(pool agent.TargetPool) error {
	for {
		fmsg, err := ft.afs.Recv()
		if err != nil {
			return ft.Error(err)
		}

		client := pool.GetTarget(fmsg.Target)
		if client == nil { //目标服务器无效,丢弃这次数据
			agent.InvalidTarget(ft, fmsg.Target)
			logger.Errorf("not find client<%d>", fmsg.Target)
		} else {
			err = agent.T2T(ft, client, fmsg)
			if err != nil { //转发失败
				logger.Error(err)
			}
		}
	}
}
