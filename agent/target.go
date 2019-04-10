package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"lemna/logger"

	context "golang.org/x/net/context"
)

type Stream interface {
	Send(*rpc.ForwardMsg) error
	Recv() (*rpc.ForwardMsg, error)
	Context() context.Context
}

type TargetPool interface {
	GetTarget(int32, *Target) *Target
	SetTargetPool(TargetPool)
}

type Target struct {
	stream Stream
	id     int32
	cache  map[int32]*Target
}

func NewTarget(s Stream, id int32) *Target {
	return &Target{stream: s, id: id, cache: make(map[int32]*Target)}
}

func (t *Target) Error(err interface{}) error {
	return fmt.Errorf("<id=%d>%s", t.id, err)
}

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
