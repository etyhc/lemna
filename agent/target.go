package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"lemna/logger"
	"sync"

	context "golang.org/x/net/context"
)

// Stream 代理目标流,接口与grpc Stream保持一致
type Stream interface {
	Send(*rpc.ForwardMsg) error
	Recv() (*rpc.ForwardMsg, error)
	Context() context.Context
}

// TargetPool 转发目标池，代理从池中得转发目标
type TargetPool interface {
	//GetTarget 得到转发目标,没有返回nil
	//    int32 目标id
	//  *Target 原目标，根据业务不同，可能需要
	GetTarget(int32, *Target) *Target
	//SetTargetPool 设置转发目标池
	SetTargetPool(TargetPool)
}

// Target 代理服务的对象统称目标，目标可以相互转发消息
type Target struct {
	stream Stream            //目标网络流
	id     int32             //目标id
	cache  map[int32]*Target //转发目标缓存
	mu     sync.Mutex
	Value  interface{} //使用者可以保存任何数据
}

// NewTarget 新代理目标
//         s 目标的网络流
//        id 目标标识，客户端唯一，服务器可能不唯一
func NewTarget(s Stream, id int32) *Target {
	return &Target{stream: s, id: id, cache: make(map[int32]*Target)}
}

// Error 附加目标信息到错误上
func (t *Target) Error(err interface{}) error {
	return fmt.Errorf("<id=%d>%s", t.id, err)
}

// Run  运行转发功能，循环等待消息并转发
// pool 转发目标池
//      等待消息错误返回
//      在自己的缓存未找到转发目标，再从转发目标池寻找转发目标，无视无转发目标错误
//      将自己缓存到转发目标
//      转发失败清除自己缓存的转发目标
func (t *Target) Run(pool TargetPool) error {
	for {
		fmsg, err := t.stream.Recv()
		if err != nil {
			return t.Error(err)
		}

		tt, ok := t.cache[fmsg.Target]
		if !ok {
			tt = pool.GetTarget(fmsg.Target, t)
			if tt == nil {
				logger.Errorf("not find target<%d>", fmsg.Target)
				continue
			}
		}
		tt.cache[t.id] = t

		//转发指令
		fmsg.Target = t.id
		err = tt.stream.Send(fmsg)
		//转发失败
		if err != nil {
			logger.Error(tt.Error(err))
			delete(t.cache, tt.id)
			continue
		}
	}
}
