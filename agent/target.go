package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"lemna/logger"

	context "golang.org/x/net/context"
)

// Stream 目标网络流
type Stream interface {
	Send(*rpc.ForwardMsg) error
	Recv() (*rpc.ForwardMsg, error)
	Context() context.Context
}

// TargetPool 代理目标池，代理从目标池得转发目标
type TargetPool interface {
	GetTarget(int32, *Target) *Target
	SetTargetPool(TargetPool)
}

// Target 代理目标，目标可以相互转发消息
type Target struct {
	stream Stream
	id     int32
	cache  map[int32]*Target
}

// NewTarget 新代理目标
// s 目标的网络流
// id 目标的id，客户端唯一，服务器可能不唯一
func NewTarget(s Stream, id int32) *Target {
	return &Target{stream: s, id: id, cache: make(map[int32]*Target)}
}

// Error 附加目标信息到错误上
func (t *Target) Error(err interface{}) error {
	return fmt.Errorf("<id=%d>%s", t.id, err)
}

// Run 运行转发功能
//     阻塞到原目标错误发生，无视转发目标错误
//     转发成功会缓存转发目标
//     缓存未找到转发目标，再从目标池寻找目标
func (t *Target) Run(pool TargetPool) error {
	for {
		fmsg, err := t.stream.Recv()
		if err != nil {
			return t.Error(err)
		}

		tt := t.cache[fmsg.Target]
		if tt == nil {
			tt = pool.GetTarget(fmsg.Target, t)
		}
		if tt == nil {
			logger.Errorf("not find target<%d>", fmsg.Target)
			continue
		}

		//转发指令
		fmsg.Target = t.id
		err = tt.stream.Send(fmsg)
		//转发失败
		if err != nil {
			logger.Error(tt.Error(err))
			delete(t.cache, tt.id)
			continue
		}
		//转发成功缓存目标
		t.cache[tt.id] = tt
	}
}
