package rpc

import (
	"fmt"
	"hash/fnv"
	"reflect"

	proto "github.com/golang/protobuf/proto"
)

/*MsgHandler 是个消息回调函数，需要实现，并注册到MsgCenter
  target 消息来源
  msg 消息本体*/
type MsgHandler func(target int32, msg interface{})

type msgInfo struct {
	elem    reflect.Type
	handler MsgHandler
}

/*MsgCenter 消息中心,封装消息工具*/
type MsgCenter struct {
	hash map[string]int32
	info map[int32]msgInfo
}

/*NewMsgCenter 新消息中心*/
func NewMsgCenter() *MsgCenter {
	return &MsgCenter{hash: make(map[string]int32), info: make(map[int32]msgInfo)}
}

/*Reg 消息注册
  msg 消息定义
  handler 消息处理函数*/
func (mc *MsgCenter) Reg(msg interface{}, handler MsgHandler) {
	name := reflect.TypeOf(msg).Elem().Name()
	h := fnv.New32a()
	_, err := h.Write([]byte(name))
	if err != nil {
		panic(name + "\n" + err.Error())
	}
	hh := int32(h.Sum32())
	if info, ok := mc.info[hh]; ok {
		panic(fmt.Errorf("Hash(%d) conflict %s %s", hh, name, info.elem.Name()))
	}
	mc.info[hh] = msgInfo{reflect.TypeOf(msg).Elem(), handler}
	mc.hash[name] = hh
}

func WrapFMNoCheck(target int32, msg proto.Message) (*ForwardMsg, error) {
	name := reflect.TypeOf(msg).Elem().Name()
	h := fnv.New32a()
	_, err := h.Write([]byte(name))
	if err != nil {
		return nil, err
	}
	hash := int32(h.Sum32())
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &ForwardMsg{Target: target, Msg: &RawMsg{Type: hash, Raw: buf}}, nil
}

func (mc *MsgCenter) wrapRaw(msg proto.Message) (*RawMsg, error) {
	hash, ok := mc.hash[reflect.TypeOf(msg).Elem().Name()]
	if !ok {
		return nil, fmt.Errorf("%s don't register", reflect.TypeOf(msg).Elem().Name())
	}
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &RawMsg{Type: hash, Raw: buf}, nil
}

/*Wrap 将消息封装为转发消息
  target 转发目标
  msg 被编码的消息*/
func (mc *MsgCenter) WrapFM(target int32, msg proto.Message) (*ForwardMsg, error) {
	raw, err := mc.wrapRaw(msg)
	if err == nil {
		return &ForwardMsg{Target: target, Msg: raw}, nil
	}
	return nil, err
}

/*WrapBroadcast 将消息封装为转发消息
  targets 转发目标数组
  msg 被编码的消息*/
func (mc *MsgCenter) WrapBM(targets []int32, msg proto.Message) (*BroadcastMsg, error) {
	raw, err := mc.wrapRaw(msg)
	if err == nil {
		return &BroadcastMsg{Targets: targets, Msg: raw}, nil
	}
	return nil, err
}

/*Handle 转发消息处理函数
  fmsg 转发消息*/
func (mc *MsgCenter) Handle(fmsg *ForwardMsg) error {
	info, ok := mc.info[fmsg.Msg.Type]
	if !ok {
		return fmt.Errorf("unregistered message type %d", fmsg.Msg.Type)
	}
	msg := reflect.New(info.elem).Interface()
	err := proto.Unmarshal(fmsg.Msg.Raw, msg.(proto.Message))
	if err != nil {
		return err
	}
	info.handler(fmsg.Target, msg)
	return nil
}
