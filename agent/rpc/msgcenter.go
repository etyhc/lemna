package rpc

import (
	"fmt"
	"hash/fnv"
	"reflect"

	"github.com/protobuf/proto"
)

/*MsgHandler 是个消息回调函数，需要实现，并注册到MsgCenter
  target 消息来源
  msg 消息本体
  stream grpc流，用于回复*/
type MsgHandler func(target int32, msg interface{}, stream interface{})

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
func (mi *MsgCenter) Reg(msg interface{}, handler MsgHandler) {
	name := reflect.TypeOf(msg.(proto.Message)).Elem().Name()
	h := fnv.New32a()
	_, err := h.Write([]byte(name))
	if err != nil {
		panic(fmt.Errorf("fnv Hash failed %s", name))
	}
	hh := int32(h.Sum32())
	if info, ok := mi.info[hh]; ok {
		panic(fmt.Errorf("Hash(%d) conflict %s %s", hh, name, info.elem.Name()))
	}
	mi.info[hh] = msgInfo{reflect.TypeOf(msg.(proto.Message)).Elem(), handler}
	mi.hash[name] = hh
}

/*Wrap 将消息封装为转发消息
  target 转发目标
  msg 被编码的消息*/
func (mi *MsgCenter) Wrap(target int32, msg proto.Message) (*ForwardMsg, error) {
	hash, ok := mi.hash[reflect.TypeOf(msg.(proto.Message)).Elem().Name()]
	if !ok {
		return nil, fmt.Errorf("%s don't register", reflect.TypeOf(msg.(proto.Message)).Elem().Name())
	}
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &ForwardMsg{Target: target, Msg: &RawMsg{Type: hash, Raw: buf}}, nil
}

/*Handle 转发消息处理函数
  fmsg 转发消息
  stream 消息流*/
func (mi *MsgCenter) Handle(fmsg *ForwardMsg, stream interface{}) error {
	info, ok := mi.info[fmsg.Msg.Type]
	if !ok {
		return fmt.Errorf("invalid message type %d", fmsg.Msg.Type)
	}
	msg := reflect.New(info.elem).Interface()
	err := proto.Unmarshal(fmsg.Msg.Raw, msg.(proto.Message))
	if err != nil {
		return err
	}
	info.handler(fmsg.Target, msg, stream)
	return nil
}
