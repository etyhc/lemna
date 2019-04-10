package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"lemna/logger"

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
	dirty  chan int32        //清理缓存chan
}

// NewTarget 新代理目标
//         s 目标的网络流
//        id 目标标识，客户端唯一，服务器可能不唯一
func NewTarget(s Stream, id int32) *Target {
	return &Target{stream: s, id: id, cache: make(map[int32]*Target), dirty: make(chan int32)}
}

// Error 附加目标信息到错误上
func (t *Target) Error(err interface{}) error {
	return fmt.Errorf("<id=%d>%s", t.id, err)
}

// Dirty 标记转发目标缓存失效
//    id 使某个转发目标缓存失效，0表示整个缓存失效
func (t *Target) Dirty(id int32) {
	t.dirty <- id
}

// Run 运行转发功能，等待消息并转发
//     等待消息错误返回，无视转发目标错误
//     转发成功会缓存转发目标
//     缓存未找到转发目标，再从转发目标池寻找转发目标
func (t *Target) Run(pool TargetPool) error {
	for {
		fmsg, err := t.stream.Recv()
		if err != nil {
			return t.Error(err)
		}
		select {
		case id := <-t.dirty:
			if id == 0 {
				t.cache = make(map[int32]*Target)
			} else {
				delete(t.cache, id)
			}
		default:
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
