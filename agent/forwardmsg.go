package agent

import (
	"fmt"
	"hash/fnv"
	"lemna/agent/rpc"
	"reflect"

	"github.com/protobuf/proto"
)

/*MsgHandler 是个消息回调函数，需要实现，并注册
  id 消息来源
  msg 消息本体
  stream grpc流，用于回复*/
type MsgHandler func(id int32, msg interface{}, stream interface{})

type msgInfo struct {
	elem    reflect.Type
	handler MsgHandler
}

var msgMap = make(map[int32]msgInfo)

/*MsgReg 消息注册
  typeid 消息类型唯一ID
  msg 消息定义
  handler 消息处理函数*/
func MsgReg(msg interface{}, handler MsgHandler) {
	elem := reflect.TypeOf(msg.(proto.Message)).Elem()
	h := fnv.New32a()
	h.Write([]byte(elem.Name()))
	hh := int32(h.Sum32())
	if _, ok := msgMap[hh]; ok {
		panic(fmt.Errorf("消息%s Hash冲突=%d", elem.Name(), hh))
	} else {
		msgMap[hh] = msgInfo{reflect.TypeOf(msg.(proto.Message)).Elem(), handler}
	}
}

/*ForwardMsgWrap 将消息封装为转发消息
  target 转发目标
  msgid 消息类型唯一id
  msg 被编码的消息*/
func ForwardMsgWrap(target int32, msg proto.Message) (*rpc.ForwardMsg, error) {
	buf, err := proto.Marshal(msg)
	h := fnv.New32a()
	h.Write([]byte(reflect.TypeOf(msg.(proto.Message)).Elem().Name()))
	return &rpc.ForwardMsg{Target: target, Msg: &rpc.RawMsg{Type: int32(h.Sum32()), Raw: buf}}, err
}

/*ForwardMsgHandle 转发消息处理函数
  id 消息来源
  typeid 消息类型唯一ID
  raw 未解码消息
  stream 消息流*/
func ForwardMsgHandle(fmsg *rpc.ForwardMsg, stream interface{}) error {
	if info, ok := msgMap[fmsg.Msg.Type]; ok {
		msg := reflect.New(info.elem).Interface()
		err := proto.Unmarshal(fmsg.Msg.Raw, msg.(proto.Message))
		if err == nil {
			info.handler(fmsg.Target, msg, stream)
		}
		return err
	}
	return fmt.Errorf("invalid message type %d", fmsg.Msg.Type)
}
