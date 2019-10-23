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
	_info  Info                      //服务器信息
}

// NewMTarget 新服务器
//    client rpc客户端
//      info 订阅的服务器信息
func NewMTarget(stream arpc.Srpc_MulticastServer, info Info) *MTarget {
	return &MTarget{stream: stream, _info: info}
}

// Send 发送转发消息给服务器
func (mt *MTarget) Send(msg *arpc.ForwardMsg) error {
	return fmt.Errorf("no send func")
}

// ID Target接口实现
//    服务器类型ID
func (mt *MTarget) ID() uint32 {
	return mt._info.Type
}

//Bind 接口实现
func (mt *MTarget) Bind(agent.Target) {
}

func (mt *MTarget) info() *Info {
	return &mt._info
}

//Forward Target接口实现
func (mt *MTarget) Forward(pool agent.TargetPool) error {
	for {
		mmsg, err := mt.stream.Recv()
		if err != nil {
			return fmt.Errorf("<type=%d>%w", mt._info.Type, err)
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
