package agent

import (
	"lemna/agent/arpc"
	"lemna/logger"
)

type Target interface {
	Send(*arpc.ForwardMsg) error
	ID() int32
	Error(err interface{}) error
}

type CTarget interface {
	Target
	Recv() (*arpc.ForwardMsg, error)
	Cache(Target)
	GetCache(int32) Target
	Uncache(int32)
}
type STarget interface {
	Target
	Recv() (*arpc.BroadcastMsg, error)
}

type TargetPool interface {
	GetTarget(int32) Target
	Bind(TargetPool)
	Run() error
}

func CtoS(src CTarget, pool TargetPool) error {
	for {
		fmsg, err := src.Recv()
		if err != nil {
			return src.Error(err)
		}

		dest := src.GetCache(fmsg.Target)
		if dest == nil {
			dest = pool.GetTarget(fmsg.Target)
			if dest == nil {
				logger.Errorf("not find server<%d>", fmsg.Target)
				continue
			}
			logger.Debug("===", dest)
		}
		logger.Debug("---", dest)

		//转发指令
		fmsg.Target = src.ID()
		err = dest.Send(fmsg)
		//转发失败
		if err != nil {
			logger.Error(dest.Error(err))
			src.Uncache(dest.ID())
			continue
		}
	}
}

//Run 运行服务器,接收服务器消息，转发给客户端
func StoC(src STarget, pool TargetPool) error {
	for {
		bmsg, err := src.Recv()
		if err != nil {
			return src.Error(err)
		}

		for _, cid := range bmsg.Targets {
			c := pool.GetTarget(cid).(CTarget)
			if c == nil {
				logger.Errorf("not find client<%d>", cid)
				continue
			}

			//转发指令
			err = c.Send(&arpc.ForwardMsg{Target: src.ID(), Msg: bmsg.Msg})
			//转发失败
			if err != nil {
				logger.Error(c.Error(err))
				c.Uncache(src.ID())
			} else {
				c.Cache(src)
			}
		}
	}
}
