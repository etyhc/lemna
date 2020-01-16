package msg

import (
	"fmt"
)

// Stream 消息流，用于接收和发送消息
//        MsgProc处理消息时会将消息流代入消息回调函数，以便消息回复.
//        这个接口只有一个用处，让服务器程序写游戏客户端测试代码能够复用msgproc代码
type Stream interface {
	// Forward 向服务器或者客户端转发消息
	Send(uint32, interface{}) error
}

// Handler 是个消息回调函数，需要实现，并注册到MsgCenter
//      int32 消息来源
//  interface 消息本体
//     Stream 消息来源流，用于回复
type Handler func(uint32, interface{}, Stream)

//Info 消息基本信息
type Info interface {
	ID() uint32   //消息ID，用于区分不同消息，客户端和服务器协商确定
	Name() string //消息名字
}

//Helper 消息辅助接口，不同消息协议自行实现
type Helper interface {
	//Extract 提取消息基本信息
	Extract(interface{}) Info
	//ToRaw 序列化消息
	ToRaw(interface{}) (uint32, []byte, error)
	//FromRaw 反序列化消息
	FromRaw(Info, []byte) (interface{}, error)
	Wrap(uint32, interface{}) (ForwardMsg, error)
	WrapMM([]uint32, interface{}) (MulticastMsg, error)
}

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
