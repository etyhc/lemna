package srpc

import (
	"fmt"
)

type msgInfo struct {
	info    MsgInfo
	handler Handler
}

/*Processor 消息处理器,封装消息工具*/
type Processor struct {
	helper Helper
	infos  map[uint32]msgInfo
}

/*NewProcessor 新消息处理器*/
func NewProcessor(helper Helper) *Processor {
	return &Processor{
		helper: helper,
		infos:  make(map[uint32]msgInfo)}
}

// Reg 消息注册
//     msg 被注册的消息
// handler 消息处理函数
func (mp *Processor) Reg(msg interface{}, handler Handler) {
	nmi := msgInfo{info: mp.helper.Extract(msg), handler: handler}
	if omi, ok := mp.infos[nmi.info.ID()]; ok {
		panic(fmt.Sprintf("MsgInfo ID(%d) conflict %s %s",
			nmi.info.ID(), nmi.info.Name(), omi.info.Name()))
	}
	mp.infos[nmi.info.ID()] = nmi
}

// Handle 处理转发消息函数
//        转发消息会被解码成原始消息，并调用注册过的处理函数来处理此原始消息
//   fmsg 收到的转发消息
//   from 消息流
func (mp *Processor) Handle(cmsg *CallMsg) (*CallMsg, error) {
	info, ok := mp.infos[cmsg.GetMid()]
	if !ok {
		return nil, fmt.Errorf("unregistered msg id=%d", cmsg.GetMid())
	}
	msg, err := mp.helper.FromRaw(info.info, cmsg.GetRaw())
	if err != nil {
		return nil, err
	}
	return info.handler(msg)
}
