package server

import (
	"fmt"
	"lemna/agent"
	"lemna/arpc"
	"lemna/logger"
)

// MTarget rpc代理端(代理服务器作为rpc服务器)
type MTarget struct {
	stream arpc.Srpc_MulticastServer //Forward调用接收、发送端
	info   *Info                     //服务器信息
}

// NewMTarget 新服务器
//    client rpc客户端
//      info 订阅的服务器信息
func NewMTarget(stream arpc.Srpc_MulticastServer) *MTarget {
	return &MTarget{stream: stream}
}

// Error 附加服务器信息到错误信息上
func (mt *MTarget) Error(err error) error {
	return fmt.Errorf("<type=%d>%w", mt.info.Type, err)
}

// Send 发送转发消息给服务器
func (mt *MTarget) Send(msg *arpc.ForwardMsg) error {
	return fmt.Errorf("no send func")
}

// ID 服务器类型ID
func (mt *MTarget) ID() uint32 {
	return mt.info.Type
}

func (mt *MTarget) Forward(pool agent.TargetPool) error {
	for {
		mmsg, err := mt.stream.Recv()
		if err != nil {
			return err
		}

		for _, cid := range mmsg.Targets {
			c := pool.GetTarget(cid)
			if c != nil {
				err = agent.T2T(mt, c, &arpc.ForwardMsg{Target: mt.ID(), Msg: mmsg.Msg})
				//转发失败
				if err != nil {
					logger.Error(err)
				}
			} else {
				logger.Errorf("not find client<%d>", cid)
			}
		}
	}
}
