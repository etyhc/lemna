package msg

import (
	"fmt"
)

type msgInfo struct {
	info    Info
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

func (mp *Processor) infoByID(id uint32) (msgInfo, error) {
	info, ok := mp.infos[id]
	if ok {
		return info, nil
	}
	return info, fmt.Errorf("unregistered msg id=%d", id)
}

func (mp *Processor) infoByMsg(msg interface{}) (msgInfo, error) {
	return mp.infoByID(mp.helper.Extract(msg).ID())
}

//UnregFM 将未注册消息封装为转发消息
func (mp *Processor) UnregFM(target uint32, msg interface{}) (ForwardMsg, error) {
	return mp.helper.Wrap(target, msg)
}

// WrapFM 将已注册消息封装为转发消息
// target 转发目标
//    msg 被转发的消息
func (mp *Processor) WrapFM(target uint32, msg interface{}) (ForwardMsg, error) {
	_, err := mp.infoByMsg(msg)
	if err != nil {
		return nil, err
	}
	return mp.helper.Wrap(target, msg)
}

// WrapMM 将已注册消息封装为指定目标的多播消息
//  targets 转发目标切片
//      msg 被广播的消息
func (mp *Processor) WrapMM(targets []uint32, msg interface{}) (MulticastMsg, error) {
	_, err := mp.infoByMsg(msg)
	if err != nil {
		return nil, err
	}
	return mp.helper.WrapMM(targets, msg)
}

// Handle 处理转发消息函数
//        转发消息会被解码成原始消息，并调用注册过的处理函数来处理此原始消息
//   fmsg 收到的转发消息
//   from 消息流
func (mp *Processor) Handle(fmsg ForwardMsg, from Stream) error {
	info, err := mp.infoByID(fmsg.GetMid())
	if err != nil {
		return err
	}
	msg, err := mp.helper.FromRaw(info.info, fmsg.GetRaw())
	if err != nil {
		return err
	}
	info.handler(fmsg.GetTarget(), msg, from)
	return nil
}
