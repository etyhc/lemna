package server

import (
	"fmt"
	"lemna/agent"
	"lemna/arpc"
	"lemna/logger"
)

// FTarget rpc代理端(代理服务器作为rpc服务器)
type FTarget struct {
	stream arpc.SAgent_ForwardServer //Forward调用接收、发送端
	_info  Info                      //服务器信息
}

// NewFTarget 新服务器
//    client rpc客户端
//      info 订阅的服务器信息
func NewFTarget(stream arpc.SAgent_ForwardServer, info Info) *FTarget {
	return &FTarget{stream: stream, _info: info}
}

func (ft *FTarget) wrapErr(err error) error {
	if err != nil {
		return fmt.Errorf("<type=%d>%w", ft._info.Type, err)
	}
	return nil
}

// Send 发送转发消息给服务器
func (ft *FTarget) Send(msg *arpc.ForwardMsg) error {
	return ft.wrapErr(ft.stream.Send(msg))
}

// ID 服务器类型ID
func (ft *FTarget) ID() uint32 {
	return ft._info.Type
}

//Bind 接口实现
func (ft *FTarget) Bind(agent.Target) {
}

func (ft *FTarget) info() *Info {
	return &ft._info
}

//Forward agent.Target接口实现
func (ft *FTarget) Forward(pool agent.TargetPool) error {
	for {
		fmsg, err := ft.stream.Recv()
		if err != nil {
			return ft.wrapErr(err)
		}

		client := pool.GetTarget(fmsg.Target)
		if client == nil { //目标服务器无效,丢弃这次数据
			agent.InvalidTarget(ft, fmsg.Target)
			logger.Errorf("not find client<%d>", fmsg.Target)
		} else {
			err = agent.T2T(ft, client, fmsg)
			if err != nil { //转发失败
				logger.Error(err)
			} else {
				//刷新客户端服务器缓冲
				client.Bind(ft)
			}
		}
	}
}
